package client

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"log/slog"
	"net/http"
	rsPb "team_dynamics/api/proto/rating_service"
	usPb "team_dynamics/api/proto/user_service"
	"team_dynamics/user_event/connection"
	"team_dynamics/user_event/downstream"
	"team_dynamics/user_event/models"
	"team_dynamics/user_event/redis"
	"strconv"
	"time"
)

type Factory interface {
	MakeClient(ctx context.Context, w http.ResponseWriter, r *http.Request) (Client, error)
}

type clientFactoryImpl struct {
	usFactory    downstream.UserServiceClientFactory
	rsFactory    downstream.RatingServiceClientFactory
	msFactory    downstream.MatchServiceClientFactory
	userKvRepo   redis.UserKvRepo
	disconnectCh chan<- Client
	logger       *slog.Logger
	config       *UserEventConfig
	upgrader     *websocket.Upgrader
}

func NewClientFactory(
	usFactory downstream.UserServiceClientFactory,
	rsFactory downstream.RatingServiceClientFactory,
	msFactory downstream.MatchServiceClientFactory,
	userKvRepo redis.UserKvRepo,
	disconnectCh chan<- Client,
	logger *slog.Logger,
	config *UserEventConfig,
	upgrader *websocket.Upgrader,
) Factory {
	return &clientFactoryImpl{
		usFactory,
		rsFactory,
		msFactory,
		userKvRepo,
		disconnectCh,
		logger,
		config,
		upgrader,
	}
}

func (f *clientFactoryImpl) getPlayerData(ctx context.Context, playerID int64) (*models.PlayerUserData, error) {
	type nameResult struct {
		name string
		err  error
	}
	type ratingResult struct {
		rating int64
		err    error
	}

	nameCh := make(chan nameResult, 1)
	ratingCh := make(chan ratingResult, 1)

	go func() {
		resp, err := f.usFactory.GetUserData(ctx, &usPb.GetUserDataRequest{Id: &playerID})
		if err != nil {
			nameCh <- nameResult{err: err}
			return
		}
		if resp.UserData == nil || resp.UserData.Name == nil {
			nameCh <- nameResult{err: fmt.Errorf("user data or name missing for player %d", playerID)}
			return
		}
		nameCh <- nameResult{name: *resp.UserData.Name}
	}()

	go func() {
		resp, err := f.rsFactory.GetRating(ctx, &rsPb.GetRatingRequest{UserId: &playerID})
		if err != nil {
			ratingCh <- ratingResult{err: err}
			return
		}
		if resp.Rating == nil || resp.Rating.DisplayValue == nil {
			ratingCh <- ratingResult{err: fmt.Errorf("rating data missing for player %d", playerID)}
			return
		}
		ratingCh <- ratingResult{rating: *resp.Rating.DisplayValue}
	}()

	nr := <-nameCh
	rr := <-ratingCh

	if nr.err != nil {
		return nil, fmt.Errorf("failed to get user name: %w", nr.err)
	}
	if rr.err != nil {
		return nil, fmt.Errorf("failed to get user rating: %w", rr.err)
	}

	return &models.PlayerUserData{
		Id:     playerID,
		Name:   nr.name,
		Rating: rr.rating,
	}, nil
}

func (f *clientFactoryImpl) validateRequest(r *http.Request) (int64, error) {
	playerIDStr := r.URL.Query().Get("playerId")
	if playerIDStr == "" {
		return 0, fmt.Errorf("missing playerId query parameter")
	}
	playerID, err := strconv.ParseInt(playerIDStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid playerId: %w", err)
	}
	return playerID, nil
}

func (f *clientFactoryImpl) MakeClient(ctx context.Context, w http.ResponseWriter, r *http.Request) (Client, error) {
	// 1. Validate request
	playerID, err := f.validateRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil, err
	}

	// 2. Get player data
	player, err := f.getPlayerData(ctx, playerID)
	if err != nil {
		f.logger.Error("failed to fetch player data", "player_id", playerID, "error", err)
		http.Error(w, "failed to fetch player data", http.StatusNotFound)
		return nil, err
	}

	// 3. Register player
	ok, err := f.userKvRepo.Register(ctx, player)
	if err != nil {
		f.logger.Error("failed to register player", "player_id", playerID, "error", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return nil, err
	}
	if !ok {
		f.logger.Warn("player already registered", "player_id", playerID)
		http.Error(w, "player already connected", http.StatusConflict)
		return nil, fmt.Errorf("player already connected")
	}

	// 4. Upgrade connection
	conn, err := f.upgrader.Upgrade(w, r, nil)
	if err != nil {
		f.logger.Error("failed to upgrade connection", "error", err)
		return nil, err
	}

	return &clientImpl{
		conn:              connection.WrapConnection(conn),
		disconnectCh:      f.disconnectCh,
		closeCh:           make(chan struct{}, 1),
		checkStatusTicker: time.NewTicker(f.config.CheckStatusPeriod),
		userKvRepo:        f.userKvRepo,
		logger:            f.logger,
		player:            player,
		msFactory:         f.msFactory,
		config:            f.config,
		lastKnownStatuses: make(map[int64]models.PlayerStatus),
	}, nil
}
