package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	pb "team_dynamics/api/proto/match_history_service"
	"team_dynamics/match_history_service/controllers"
	"team_dynamics/match_history_service/pg"
	"team_dynamics/match_history_service/services"
	pglib "team_dynamics/pg_lib/include"
	"time"
)

func getPgConnectionConfig() (*pglib.ConnectionConfig, error) {
	host := os.Getenv("PG_HOST")
	if host == "" {
		return nil, errors.New("PG_HOST environment variable not set")
	}

	portStr := os.Getenv("PG_PORT")
	if portStr == "" {
		return nil, errors.New("PG_PORT environment variable not set")
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid PG_PORT value %q: %w", portStr, err)
	}

	user := os.Getenv("PG_USER")
	if user == "" {
		return nil, errors.New("PG_USER environment variable not set")
	}

	password := os.Getenv("PG_PASSWORD")
	if password == "" {
		return nil, errors.New("PG_PASSWORD environment variable not set")
	}

	database := os.Getenv("PG_DATABASE")
	if database == "" {
		return nil, errors.New("PG_DATABASE environment variable not set")
	}

	sslMode := os.Getenv("PG_SSLMODE")
	if sslMode == "" {
		sslMode = "disable"
	}

	return &pglib.ConnectionConfig{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		Database: database,
		SSLMode:  sslMode,
	}, nil
}

func getPgInitConfig() (*pglib.InitializationConfig, error) {
	timeoutStr := os.Getenv("PG_CONN_TIMEOUT")
	if timeoutStr == "" {
		return nil, errors.New("PG_CONN_TIMEOUT environment variable not set")
	}
	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		return nil, fmt.Errorf("invalid PG_CONN_TIMEOUT value %q: %w", timeoutStr, err)
	}

	retriesStr := os.Getenv("PG_RETRIES")
	if retriesStr == "" {
		return nil, errors.New("PG_RETRIES environment variable not set")
	}
	retries, err := strconv.ParseInt(retriesStr, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid PG_RETRIES value %q: %w", retriesStr, err)
	}

	return &pglib.InitializationConfig{
		ConnectionTimeout: timeout,
		Retries:           int32(retries),
	}, nil
}

func getPgPoolConfig() (*pglib.PoolConfig, error) {
	maxConnsStr := os.Getenv("PG_MAX_CONNS")
	if maxConnsStr == "" {
		return nil, errors.New("PG_MAX_CONNS environment variable not set")
	}
	maxConns, err := strconv.ParseInt(maxConnsStr, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid PG_MAX_CONNS value %q: %w", maxConnsStr, err)
	}

	minConnsStr := os.Getenv("PG_MIN_CONNS")
	if minConnsStr == "" {
		return nil, errors.New("PG_MIN_CONNS environment variable not set")
	}
	minConns, err := strconv.ParseInt(minConnsStr, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid PG_MIN_CONNS value %q: %w", minConnsStr, err)
	}

	lifetimeStr := os.Getenv("PG_MAX_CONN_LIFETIME")
	if lifetimeStr == "" {
		return nil, errors.New("PG_MAX_CONN_LIFETIME environment variable not set")
	}
	lifetime, err := time.ParseDuration(lifetimeStr)
	if err != nil {
		return nil, fmt.Errorf("invalid PG_MAX_CONN_LIFETIME value %q: %w", lifetimeStr, err)
	}

	idleTimeStr := os.Getenv("PG_MAX_CONN_IDLE_TIME")
	if idleTimeStr == "" {
		return nil, errors.New("PG_MAX_CONN_IDLE_TIME environment variable not set")
	}
	idleTime, err := time.ParseDuration(idleTimeStr)
	if err != nil {
		return nil, fmt.Errorf("invalid PG_MAX_CONN_IDLE_TIME value %q: %w", idleTimeStr, err)
	}

	healthCheckStr := os.Getenv("PG_HEALTH_CHECK_PERIOD")
	if healthCheckStr == "" {
		return nil, errors.New("PG_HEALTH_CHECK_PERIOD environment variable not set")
	}
	healthCheck, err := time.ParseDuration(healthCheckStr)
	if err != nil {
		return nil, fmt.Errorf("invalid PG_HEALTH_CHECK_PERIOD value %q: %w", healthCheckStr, err)
	}

	return &pglib.PoolConfig{
		MaxConns:          int32(maxConns),
		MinConns:          int32(minConns),
		MaxConnLifetime:   lifetime,
		MaxConnIdleTime:   idleTime,
		HealthCheckPeriod: healthCheck,
	}, nil
}

func getListenAddress() (string, error) {
	listenAddress := os.Getenv("LISTEN_ADDRESS")
	if listenAddress == "" {
		return "", fmt.Errorf("LISTEN_ADDRESS environment vaiable not set")
	}
	return listenAddress, nil
}

func getHttpListenAddress() (string, error) {
	listenAddress := os.Getenv("HTTP_LISTEN_ADDRESS")
	if listenAddress == "" {
		return "", fmt.Errorf("HTTP_LISTEN_ADDRESS environment vaiable not set")
	}
	return listenAddress, nil
}

func main() {
	pgConnCfg, err := getPgConnectionConfig()
	if err != nil {
		panic(fmt.Sprintf("Unable to get pg connection config: %v", err))
	}
	pgInitCfg, err := getPgInitConfig()
	if err != nil {
		panic(fmt.Sprintf("Unable to get pg init config: %v", err))
	}
	pgPoolCfg, err := getPgPoolConfig()
	if err != nil {
		panic(fmt.Sprintf("Unable to get pg pool config: %v", err))
	}

	pool, err := pglib.MakePool(context.Background(), pgConnCfg, pgPoolCfg, pgInitCfg)
	if err != nil {
		panic(fmt.Sprintf("Unable to get pg pool: %v", err))
	}

	listenAddress, err := getListenAddress()
	if err != nil {
		panic(fmt.Sprintf("Unable to get listen address: %v", err))
	}

	httpListenAddress, err := getHttpListenAddress()
	if err != nil {
		panic(fmt.Sprintf("Unable to get http listen address: %v", err))
	}

	controller := &controllers.MatchHistoryServiceController{
		Service: services.MakeMatchHistoryService(
			services.MakePageKeyService(),
			pg.MakeMatchHistoryRepo(pool),
		),
	}

	lis, err := net.Listen("tcp", listenAddress)
	if err != nil {
		panic(fmt.Sprintf("Unable to listen: %v", err))
	}

	s := grpc.NewServer()
	pb.RegisterMatchHistoryServiceServer(s, controller)

	go func() {
		if err := s.Serve(lis); err != nil {
			panic(fmt.Sprintf("Unable to serve: %v", err))
		}
	}()
	ctx := context.Background()
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err = pb.RegisterMatchHistoryServiceHandlerFromEndpoint(ctx, mux, listenAddress, opts)
	if err != nil {
		log.Fatalf("error while running grpc gateway, %v", err)
	}
	if err := http.ListenAndServe(httpListenAddress, mux); err != nil {
		log.Fatalf("error on listening")
	}
}
