package main

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"os"
	"strconv"
	pb "team_dynamics/api/proto/rating_service"
	pglib "team_dynamics/pg_lib/include"
	"team_dynamics/rating_service/controllers"
	"team_dynamics/rating_service/pg"
	"team_dynamics/rating_service/services"
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

	controller := &controllers.RatingServiceController{
		Service: services.MakeRatingService(pg.MakeRatingRepo(pool), services.MakeGlicko2Service()),
	}

	lis, err := net.Listen("tcp", listenAddress)
	if err != nil {
		panic(fmt.Sprintf("Unable to listen: %v", err))
	}

	s := grpc.NewServer()
	pb.RegisterRatingServiceServer(s, controller)

	if err := s.Serve(lis); err != nil {
		panic(fmt.Sprintf("Unable to serve: %v", err))
	}
}
