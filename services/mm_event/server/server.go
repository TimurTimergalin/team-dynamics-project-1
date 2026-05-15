package main

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"team_dynamics/mm_event/client"
	"team_dynamics/mm_event/controllers"
	"team_dynamics/mm_event/downstream"
	mmeRedis "team_dynamics/mm_event/redis"
	"time"
)

type MMEventConfig struct {
	ListenAddress               string
	MatchServiceAddress         string
	MatchHistoryServiceV2Address string
	UserServiceAddress          string
	ChannelSizes                int
}

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

func getMMEventConfig() (*MMEventConfig, error) {
	listenAddr := os.Getenv("LISTEN_ADDRESS")
	if listenAddr == "" {
		return nil, errors.New("LISTEN_ADDRESS environment variable not set")
	}
	matchAddr := os.Getenv("MATCH_SERVICE_ADDRESS")
	if matchAddr == "" {
		return nil, errors.New("MATCH_SERVICE_ADDRESS environment variable not set")
	}
	v2Addr := os.Getenv("MATCH_HISTORY_SERVICE_V2_ADDRESS")
	if v2Addr == "" {
		return nil, errors.New("MATCH_HISTORY_SERVICE_V2_ADDRESS environment variable not set")
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

	return &MMEventConfig{
		ListenAddress:                listenAddr,
		MatchServiceAddress:          matchAddr,
		MatchHistoryServiceV2Address: v2Addr,
		UserServiceAddress:           userAddr,
		ChannelSizes:                 channelSizes,
	}, nil
}

func getClientConfig() (*client.Config, error) {
	msgTimeoutStr := os.Getenv("MESSAGE_RECEIVED_TIMEOUT")
	if msgTimeoutStr == "" {
		return nil, errors.New("MESSAGE_RECEIVED_TIMEOUT environment variable not set")
	}
	msgTimeout, err := time.ParseDuration(msgTimeoutStr)
	if err != nil {
		return nil, err
	}

	checkMatchStr := os.Getenv("CHECK_MATCH_PERIOD")
	if checkMatchStr == "" {
		return nil, errors.New("CHECK_MATCH_PERIOD environment variable not set")
	}
	checkMatch, err := time.ParseDuration(checkMatchStr)
	if err != nil {
		return nil, err
	}

	checkInPoolStr := os.Getenv("CHECK_IN_POOL_PERIOD")
	if checkInPoolStr == "" {
		return nil, errors.New("CHECK_IN_POOL_PERIOD environment variable not set")
	}
	checkInPool, err := time.ParseDuration(checkInPoolStr)
	if err != nil {
		return nil, err
	}

	hubTimeoutStr := os.Getenv("HUB_REGISTER_TIMEOUT")
	if hubTimeoutStr == "" {
		return nil, errors.New("HUB_REGISTER_TIMEOUT environment variable not set")
	}
	hubTimeout, err := time.ParseDuration(hubTimeoutStr)
	if err != nil {
		return nil, err
	}

	connectionTtlString := os.Getenv("CONNECTION_TTL")
	if connectionTtlString == "" {
		return nil, errors.New("CONNECTION_TTL env is not set")
	}
	connectionTtl, err := time.ParseDuration(connectionTtlString)
	if err != nil {
		return nil, err
	}

	return &client.Config{
		MessageReceivedTimeout: msgTimeout,
		CheckMatchPeriod:       checkMatch,
		CheckInPoolPeriod:      checkInPool,
		HubRegisterTimeout:     hubTimeout,
		ConnectionTTL:          connectionTtl,
	}, nil
}

func getMMPoolRepoConfig() (*mmeRedis.MMPoolRepoConfig, error) {
	lockTTLStr := os.Getenv("MM_LOCK_TTL")
	if lockTTLStr == "" {
		return nil, errors.New("MM_LOCK_TTL environment variable not set")
	}
	lockTTL, err := time.ParseDuration(lockTTLStr)
	if err != nil {
		return nil, fmt.Errorf("invalid MM_LOCK_TTL: %w", err)
	}
	return &mmeRedis.MMPoolRepoConfig{
		LockTTL: lockTTL,
	}, nil
}

func main() {
	mmeCfg, err := getMMEventConfig()
	if err != nil {
		log.Fatalf("cannot get mmeconfig: %v", err)
	}
	clientConfig, err := getClientConfig()
	if err != nil {
		log.Fatalf("cannot get clientconfig: %v", err)
	}
	redisOpts, err := getRedisOptions()
	if err != nil {
		log.Fatalf("cannot get redis options: %v", err)
	}
	mmPoolConfig, err := getMMPoolRepoConfig()
	if err != nil {
		log.Fatalf("cannot get mmpool configL %v", err)
	}
	redisClient := redis.NewClient(redisOpts)
	mmPoolRepo := mmeRedis.MakeMMPoolRepo(redisClient, mmPoolConfig)

	registerCh := make(chan client.Client, mmeCfg.ChannelSizes)
	disconnectCh := make(chan client.Client, mmeCfg.ChannelSizes)
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug})
	logger := slog.New(handler)

	upgrader := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(*http.Request) bool {
			return true
		},
	}
	hub := client.NewHub(disconnectCh, registerCh, logger, clientConfig)
	factory := client.NewClientFactory(
		downstream.NewMatchServiceClientFactory(mmeCfg.MatchServiceAddress),
		downstream.NewMatchHistoryServiceV2ClientFactory(mmeCfg.MatchHistoryServiceV2Address),
		downstream.NewUserServiceClientFactory(mmeCfg.UserServiceAddress),
		mmPoolRepo, disconnectCh, logger, clientConfig, upgrader,
	)

	controller := controllers.NewMMEventController(
		hub,
		factory,
		mmeCfg.ListenAddress,
	)
	controller.Run()
}
