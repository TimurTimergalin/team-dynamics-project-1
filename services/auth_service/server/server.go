package main

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	pb "team_dynamics/api/proto/auth_service"
	"team_dynamics/auth_service/controllers"
	"team_dynamics/auth_service/downstream"
	"team_dynamics/auth_service/models"
	"team_dynamics/auth_service/repos"
	"team_dynamics/auth_service/services"
	"time"
)

func requireEnv(key string) (string, error) {
	val := os.Getenv(key)
	if val == "" {
		return "", fmt.Errorf("%s environment variable not set", key)
	}
	return val, nil
}

func parseDurationEnv(key string) (time.Duration, error) {
	val, err := requireEnv(key)
	if err != nil {
		return 0, err
	}
	d, err := time.ParseDuration(val)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %w", key, err)
	}
	return d, nil
}

func loadKeyPair(privateKeyPath, publicKeyPath string) (models.KeyPair, error) {
	privBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return models.KeyPair{}, fmt.Errorf("failed to read private key file %s: %w", privateKeyPath, err)
	}
	privBlock, _ := pem.Decode(privBytes)
	if privBlock == nil {
		return models.KeyPair{}, fmt.Errorf("failed to decode PEM from %s", privateKeyPath)
	}
	privKey, err := x509.ParsePKCS8PrivateKey(privBlock.Bytes)
	if err != nil {
		return models.KeyPair{}, fmt.Errorf("failed to parse private key: %w", err)
	}
	rsaPrivKey, ok := privKey.(*rsa.PrivateKey)
	if !ok {
		return models.KeyPair{}, errors.New("private key is not RSA")
	}

	pubBytes, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return models.KeyPair{}, fmt.Errorf("failed to read public key file %s: %w", publicKeyPath, err)
	}
	pubBlock, _ := pem.Decode(pubBytes)
	if pubBlock == nil {
		return models.KeyPair{}, fmt.Errorf("failed to decode PEM from %s", publicKeyPath)
	}
	pubKeyAny, err := x509.ParsePKIXPublicKey(pubBlock.Bytes)
	if err != nil {
		return models.KeyPair{}, fmt.Errorf("failed to parse public key: %w", err)
	}
	rsaPubKey, ok := pubKeyAny.(*rsa.PublicKey)
	if !ok {
		return models.KeyPair{}, errors.New("public key is not RSA")
	}

	return models.KeyPair{PrivateKey: rsaPrivKey, PublicKey: rsaPubKey}, nil
}

func getRedisOptions() (*redis.Options, error) {
	addr, err := requireEnv("REDIS_ADDR")
	if err != nil {
		return nil, err
	}
	password, err := requireEnv("REDIS_PASSWORD")
	if err != nil {
		return nil, err
	}
	dbStr, err := requireEnv("REDIS_DB")
	if err != nil {
		return nil, err
	}
	db, err := strconv.Atoi(dbStr)
	if err != nil {
		return nil, fmt.Errorf("invalid REDIS_DB: %w", err)
	}
	dialTimeout, err := parseDurationEnv("REDIS_DIAL_TIMEOUT")
	if err != nil {
		return nil, err
	}
	readTimeout, err := parseDurationEnv("REDIS_READ_TIMEOUT")
	if err != nil {
		return nil, err
	}
	writeTimeout, err := parseDurationEnv("REDIS_WRITE_TIMEOUT")
	if err != nil {
		return nil, err
	}
	return &redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		DialTimeout:  dialTimeout,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}, nil
}

func main() {
	listenAddress, err := requireEnv("LISTEN_ADDRESS")
	if err != nil {
		log.Fatalf("cannot get listen address: %v", err)
	}
	httpListenAddress, err := requireEnv("HTTP_LISTEN_ADDRESS")
	if err != nil {
		log.Fatalf("cannot get http listen address: %v", err)
	}

	steamApiKey, err := requireEnv("STEAM_API_KEY")
	if err != nil {
		log.Fatalf("%v", err)
	}
	steamAppId, err := requireEnv("STEAM_APP_ID")
	if err != nil {
		log.Fatalf("%v", err)
	}
	issuer, err := requireEnv("JWT_ISSUER")
	if err != nil {
		log.Fatalf("%v", err)
	}
	userServiceAddress, err := requireEnv("USER_SERVICE_ADDRESS")
	if err != nil {
		log.Fatalf("%v", err)
	}
	accessExpiration, err := parseDurationEnv("JWT_ACCESS_EXPIRATION")
	if err != nil {
		log.Fatalf("%v", err)
	}
	refreshExpiration, err := parseDurationEnv("JWT_REFRESH_EXPIRATION")
	if err != nil {
		log.Fatalf("%v", err)
	}

	primaryPrivKeyPath, err := requireEnv("PRIMARY_PRIVATE_KEY_PATH")
	if err != nil {
		log.Fatalf("%v", err)
	}
	primaryPubKeyPath, err := requireEnv("PRIMARY_PUBLIC_KEY_PATH")
	if err != nil {
		log.Fatalf("%v", err)
	}
	secondaryPrivKeyPath, err := requireEnv("SECONDARY_PRIVATE_KEY_PATH")
	if err != nil {
		log.Fatalf("%v", err)
	}
	secondaryPubKeyPath, err := requireEnv("SECONDARY_PUBLIC_KEY_PATH")
	if err != nil {
		log.Fatalf("%v", err)
	}

	primaryKeyPair, err := loadKeyPair(primaryPrivKeyPath, primaryPubKeyPath)
	if err != nil {
		log.Fatalf("cannot load primary key pair: %v", err)
	}
	secondaryKeyPair, err := loadKeyPair(secondaryPrivKeyPath, secondaryPubKeyPath)
	if err != nil {
		log.Fatalf("cannot load secondary key pair: %v", err)
	}

	keyPairVersion, err := requireEnv("KEY_PAIR_VERSION")
	if err != nil {
		log.Fatalf("%v", err)
	}

	redisOpts, err := getRedisOptions()
	if err != nil {
		log.Fatalf("cannot get redis options: %v", err)
	}
	redisClient := redis.NewClient(redisOpts)

	jwtService := services.NewJwtService(
		accessExpiration,
		refreshExpiration,
		issuer,
		primaryKeyPair,
		secondaryKeyPair,
	)
	steamService := services.NewSteamService(steamApiKey, steamAppId)
	userServiceClient := downstream.NewUserServiceClientFactory(userServiceAddress)
	authKvRepo := repos.MakeAuthKvRepo(redisClient, refreshExpiration)

	authService, err := services.NewAuthService(jwtService, steamService, userServiceClient, authKvRepo, primaryKeyPair, secondaryKeyPair, keyPairVersion)
	if err != nil {
		log.Fatalf("cannot create auth service: %v", err)
	}

	controller := &controllers.AuthServiceController{Service: authService}

	lis, err := net.Listen("tcp", listenAddress)
	if err != nil {
		log.Fatalf("cannot listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterAuthServiceServer(s, controller)

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("cannot serve: %v", err)
		}
	}()

	ctx := context.Background()
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := pb.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, listenAddress, opts); err != nil {
		log.Fatalf("cannot register grpc gateway: %v", err)
	}

	if err := http.ListenAndServe(httpListenAddress, mux); err != nil {
		log.Fatalf("cannot serve http: %v", err)
	}
}
