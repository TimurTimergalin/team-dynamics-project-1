package main

import (
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"os"
	"strconv"
	fmPb "team_dynamics/api/proto/fleet_manager"
	mhsPb "team_dynamics/api/proto/match_history_service"
	pb "team_dynamics/api/proto/match_service"
	rsPb "team_dynamics/api/proto/rating_service"
	"team_dynamics/match_service/config"
	"team_dynamics/match_service/controllers"
	msRedis "team_dynamics/match_service/redis"
	"team_dynamics/match_service/services"
	"time"
)

func getRedisOptions() (*redis.Options, error) {
	// Required variables
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		return nil, errors.New("REDIS_ADDR environment variable not set")
	}
	password := os.Getenv("REDIS_PASSWORD")
	if password == "" {
		return nil, errors.New("REDIS_PASSWORD environment variable not set")
	}
	dbStr := os.Getenv("REDIS_DB")
	if dbStr == "" {
		return nil, errors.New("REDIS_DB environment variable not set")
	}
	db, err := strconv.Atoi(dbStr)
	if err != nil {
		return nil, fmt.Errorf("invalid REDIS_DB: %w", err)
	}

	poolSizeStr := os.Getenv("REDIS_POOL_SIZE")
	if poolSizeStr == "" {
		return nil, errors.New("REDIS_POOL_SIZE environment variable not set")
	}
	poolSize, err := strconv.Atoi(poolSizeStr)
	if err != nil {
		return nil, fmt.Errorf("invalid REDIS_POOL_SIZE: %w", err)
	}

	minIdleConnsStr := os.Getenv("REDIS_MIN_IDLE_CONNS")
	if minIdleConnsStr == "" {
		return nil, errors.New("REDIS_MIN_IDLE_CONNS environment variable not set")
	}
	minIdleConns, err := strconv.Atoi(minIdleConnsStr)
	if err != nil {
		return nil, fmt.Errorf("invalid REDIS_MIN_IDLE_CONNS: %w", err)
	}

	dialTimeoutStr := os.Getenv("REDIS_DIAL_TIMEOUT")
	if dialTimeoutStr == "" {
		return nil, errors.New("REDIS_DIAL_TIMEOUT environment variable not set")
	}
	dialTimeout, err := time.ParseDuration(dialTimeoutStr)
	if err != nil {
		return nil, fmt.Errorf("invalid REDIS_DIAL_TIMEOUT: %w", err)
	}

	readTimeoutStr := os.Getenv("REDIS_READ_TIMEOUT")
	if readTimeoutStr == "" {
		return nil, errors.New("REDIS_READ_TIMEOUT environment variable not set")
	}
	readTimeout, err := time.ParseDuration(readTimeoutStr)
	if err != nil {
		return nil, fmt.Errorf("invalid REDIS_READ_TIMEOUT: %w", err)
	}

	writeTimeoutStr := os.Getenv("REDIS_WRITE_TIMEOUT")
	if writeTimeoutStr == "" {
		return nil, errors.New("REDIS_WRITE_TIMEOUT environment variable not set")
	}
	writeTimeout, err := time.ParseDuration(writeTimeoutStr)
	if err != nil {
		return nil, fmt.Errorf("invalid REDIS_WRITE_TIMEOUT: %w", err)
	}

	poolTimeoutStr := os.Getenv("REDIS_POOL_TIMEOUT")
	if poolTimeoutStr == "" {
		return nil, errors.New("REDIS_POOL_TIMEOUT environment variable not set")
	}
	poolTimeout, err := time.ParseDuration(poolTimeoutStr)
	if err != nil {
		return nil, fmt.Errorf("invalid REDIS_POOL_TIMEOUT: %w", err)
	}

	return &redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		PoolSize:     poolSize,
		MinIdleConns: minIdleConns,
		DialTimeout:  dialTimeout,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		PoolTimeout:  poolTimeout,
	}, nil
}

func getMatchServiceConfig() (*config.MatchServiceConfig, error) {
	listenAddr := os.Getenv("MATCH_SERVICE_LISTEN_ADDRESS")
	if listenAddr == "" {
		return nil, errors.New("MATCH_SERVICE_LISTEN_ADDRESS environment variable not set")
	}
	fmAddr := os.Getenv("FLEET_MANAGER_ADDRESS")
	if fmAddr == "" {
		return nil, errors.New("FLEET_MANAGER_ADDRESS environment variable not set")
	}
	rsAddr := os.Getenv("RATING_SERVICE_ADDRESS")
	if rsAddr == "" {
		return nil, errors.New("RATING_SERVICE_ADDRESS environment variable not set")
	}
	mhsAddr := os.Getenv("MATCH_HISTORY_SERVICE_ADDRESS")
	if mhsAddr == "" {
		return nil, errors.New("MATCH_HISTORY_SERVICE_ADDRESS environment variable not set")
	}
	return &config.MatchServiceConfig{
		ListenAddress:              listenAddr,
		FleetManagerAddress:        fmAddr,
		RatingServiceAddress:       rsAddr,
		MatchHistoryServiceAddress: mhsAddr,
	}, nil
}

func main() {
	redisOptions, err := getRedisOptions()
	if err != nil {
		log.Fatalf("error in redis options: %v", err)
	}
	redisClient := redis.NewClient(redisOptions)
	matchServiceConfig, err := getMatchServiceConfig()
	if err != nil {
		log.Fatalf("error in match service config: %v", err)
	}
	fmConn, err := grpc.NewClient(matchServiceConfig.FleetManagerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("error in fleet manager: %v", err)
	}
	fmClient := fmPb.NewFleetManagerClient(fmConn)
	rsConn, err := grpc.NewClient(matchServiceConfig.RatingServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("error in rating service: %v", err)
	}
	rsClient := rsPb.NewRatingServiceClient(rsConn)
	mhsConn, err := grpc.NewClient(matchServiceConfig.MatchHistoryServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("error in match history service: %v", err)
	}
	mhsClient := mhsPb.NewMatchHistoryServiceClient(mhsConn)
	controller := &controllers.MatchServiceController{
		Service: services.MakeMatchService(
			msRedis.MakeMatchKvRepo(
				redisClient,
			),
			fmClient,
			rsClient,
			mhsClient,
		),
	}

	lis, err := net.Listen("tcp", matchServiceConfig.ListenAddress)
	if err != nil {
		log.Fatalf("cannot listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMatchServiceServer(s, controller)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error while trying to serve: %s", err.Error())
	}
}
