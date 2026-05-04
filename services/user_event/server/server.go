package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"team_dynamics/user_event/client"
	"team_dynamics/user_event/controllers"
	"team_dynamics/user_event/downstream"
	ueRedis "team_dynamics/user_event/redis"
	"time"
)

type UserEventServerConfig struct {
	ListenAddress        string
	MatchServiceAddress  string
	RatingServiceAddress string
	UserServiceAddress   string
	ChannelSizes         int
}

func getRedisOptions() (*redis.Options, error) {
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
	poolTimeout, err := parseDurationEnv("REDIS_POOL_TIMEOUT")
	if err != nil {
		return nil, err
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

func getServerConfig() (*UserEventServerConfig, error) {
	listenAddr := os.Getenv("LISTEN_ADDRESS")
	if listenAddr == "" {
		return nil, errors.New("LISTEN_ADDRESS environment variable not set")
	}
	matchAddr := os.Getenv("MATCH_SERVICE_ADDRESS")
	if matchAddr == "" {
		return nil, errors.New("MATCH_SERVICE_ADDRESS environment variable not set")
	}
	ratingAddr := os.Getenv("RATING_SERVICE_ADDRESS")
	if ratingAddr == "" {
		return nil, errors.New("RATING_SERVICE_ADDRESS environment variable not set")
	}
	userAddr := os.Getenv("USER_SERVICE_ADDRESS")
	if userAddr == "" {
		return nil, errors.New("USER_SERVICE_ADDRESS environment variable not set")
	}
	channelSizesStr := os.Getenv("CHANNEL_SIZES")
	if channelSizesStr == "" {
		return nil, errors.New("CHANNEL_SIZES environment variable not set")
	}
	channelSizes, err := strconv.Atoi(channelSizesStr)
	if err != nil {
		return nil, errors.New("invalid CHANNEL_SIZES: must be an integer")
	}
	return &UserEventServerConfig{
		ListenAddress:        listenAddr,
		MatchServiceAddress:  matchAddr,
		RatingServiceAddress: ratingAddr,
		UserServiceAddress:   userAddr,
		ChannelSizes:         channelSizes,
	}, nil
}

func getClientConfig() (*client.UserEventConfig, error) {
	checkStatus, err := parseDurationEnv("CHECK_STATUS_PERIOD")
	if err != nil {
		return nil, err
	}
	return &client.UserEventConfig{
		CheckStatusPeriod: checkStatus,
	}, nil
}

func parseDurationEnv(key string) (time.Duration, error) {
	val := os.Getenv(key)
	if val == "" {
		return 0, fmt.Errorf("%s environment variable not set", key)
	}
	d, err := time.ParseDuration(val)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %w", key, err)
	}
	return d, nil
}

func main() {
	serverCfg, err := getServerConfig()
	if err != nil {
		log.Fatalf("cannot get server config: %v", err)
	}
	clientCfg, err := getClientConfig()
	if err != nil {
		log.Fatalf("cannot get client config: %v", err)
	}
	redisOpts, err := getRedisOptions()
	if err != nil {
		log.Fatalf("cannot get redis options: %v", err)
	}

	redisClient := redis.NewClient(redisOpts)
	userKvRepo := ueRedis.MakeUserKvRepo(redisClient)

	registerCh := make(chan client.Client, serverCfg.ChannelSizes)
	disconnectCh := make(chan client.Client, serverCfg.ChannelSizes)

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug})
	logger := slog.New(handler)

	upgrader := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(*http.Request) bool {
			return true
		},
	}

	ctx := context.Background()
	hub := client.NewHub(disconnectCh, registerCh, logger, clientCfg)
	factory := client.NewClientFactory(
		downstream.NewUserServiceClientFactory(serverCfg.UserServiceAddress),
		downstream.NewRatingServiceClientFactory(serverCfg.RatingServiceAddress),
		downstream.NewMatchServiceClientFactory(serverCfg.MatchServiceAddress),
		userKvRepo, disconnectCh, logger, clientCfg, upgrader,
	)

	controller := controllers.NewUserEventController(ctx, hub, factory, serverCfg.ListenAddress)
	controller.Run(ctx)
}
