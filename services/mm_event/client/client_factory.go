package client

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log/slog"
	"net/http"
	"strconv"
	msPb "team_dynamics/api/proto/match_service"
	rsPb "team_dynamics/api/proto/rating_service"
	usPb "team_dynamics/api/proto/user_service"
	"team_dynamics/mm_event/connection"
	mmeJson "team_dynamics/mm_event/json"
	"team_dynamics/mm_event/models"
	"team_dynamics/mm_event/redis"
	"time"
)

type Factory interface {
	MakeClient(w http.ResponseWriter, r *http.Request) (Client, error)
}

type clientFactoryImpl struct {
	msClient     msPb.MatchServiceClient
	rsClient     rsPb.RatingServiceClient
	usClient     usPb.UserServiceClient
	mmPoolRepo   redis.MMPoolRepo
	disconnectCh chan<- Client
	logger       *slog.Logger
	config       *Config
	upgrader     *websocket.Upgrader
}

func NewClientFactory(
	msClient msPb.MatchServiceClient,
	rsClient rsPb.RatingServiceClient,
	usClient usPb.UserServiceClient,
	mmPoolRepo redis.MMPoolRepo,
	disconnectCh chan<- Client,
	logger *slog.Logger,
	config *Config,
	upgrader *websocket.Upgrader,
) Factory {
	return &clientFactoryImpl{
		msClient,
		rsClient,
		usClient,
		mmPoolRepo,
		disconnectCh,
		logger,
		config,
		upgrader,
	}
}

func (f *clientFactoryImpl) getUserData(ctx context.Context, playerID int64) (*string, error) {
	req := &usPb.GetUserDataRequest{
		Id: &playerID,
	}
	resp, err := f.usClient.GetUserData(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.UserData == nil || resp.UserData.Name == nil {
		return nil, fmt.Errorf("user data or name missing for player %d", playerID)
	}
	return resp.UserData.Name, nil
}
func (f *clientFactoryImpl) getUserRating(ctx context.Context, playerID int64) (float64, int64, error) {
	req := &rsPb.GetRatingRequest{
		UserId: &playerID,
	}
	resp, err := f.rsClient.GetRating(ctx, req)
	if err != nil {
		return 0, 0, err
	}
	if resp.Rating == nil {
		return 0, 0, fmt.Errorf("rating data missing for player %d", playerID)
	}
	rating := 0.0
	if resp.Rating.RatingValue != nil {
		rating = *resp.Rating.RatingValue
	}
	displayRating := int64(0)
	if resp.Rating.DisplayValue != nil {
		displayRating = *resp.Rating.DisplayValue
	}
	return rating, displayRating, nil
}

func (f *clientFactoryImpl) getPlayerData(ctx context.Context, playerID int64, fleet string) *models.Player {
	var name *string
	var rating float64
	var displayRating int64
	var nameErr, ratingErr error

	done := make(chan struct{}, 2)

	go func() {
		name, nameErr = f.getUserData(ctx, playerID)
		done <- struct{}{}
	}()
	go func() {
		rating, displayRating, ratingErr = f.getUserRating(ctx, playerID)
		done <- struct{}{}
	}()

	for i := 0; i < 2; i++ {
		<-done
	}

	if nameErr != nil || ratingErr != nil {
		f.logger.Error("failed to fetch player data", "player_id", playerID, "name_err", nameErr, "rating_err", ratingErr)
		return nil
	}

	return &models.Player{
		Id:              playerID,
		Name:            *name,
		Rating:          rating,
		DisplayedRating: displayRating,
		RegId:           uuid.New().String(),
		Fleet:           fleet,
	}
}

func (f *clientFactoryImpl) getMatchAddress(ctx context.Context, playerID int64) *string {
	req := &msPb.GetMatchRequest{
		PlayerId: &playerID,
	}
	resp, err := f.msClient.GetMatch(ctx, req)
	if err != nil {
		f.logger.Error("failed to get match address", "player_id", playerID, "error", err)
		return nil
	}
	if resp.ConnectionInfo == nil || resp.ConnectionInfo.Address == nil {
		return nil
	}
	return resp.ConnectionInfo.Address
}

func (f *clientFactoryImpl) validateRequest(r *http.Request) (*int64, *string, error) {
	playerIDStr := r.URL.Query().Get("playerId")
	if playerIDStr == "" {
		return nil, nil, fmt.Errorf("missing playerId query parameter")
	}
	playerID, err := strconv.ParseInt(playerIDStr, 10, 64)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid playerId: %w", err)
	}
	fleet := r.URL.Query().Get("fleet")
	if fleet == "" {
		return nil, nil, fmt.Errorf("missing fleet query parameter")
	}
	return &playerID, &fleet, nil
}

func (f *clientFactoryImpl) MakeClient(w http.ResponseWriter, r *http.Request) (Client, error) {
	ctx := context.Background()

	// 1. Validate request
	playerID, fleet, err := f.validateRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil, err
	}

	// 1b. Register connection
	ok, err := f.mmPoolRepo.AddConnection(ctx, *playerID, f.config.ConnectionTTL)
	if err != nil {
		f.logger.Error("failed to add connection", "player_id", *playerID, "error", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return nil, err
	}
	if !ok {
		f.logger.Warn("player already has an active connection", "player_id", *playerID)
		http.Error(w, "player already connected", http.StatusConflict)
		return nil, fmt.Errorf("player already connected")
	}

	// 2. Get player data
	player := f.getPlayerData(ctx, *playerID, *fleet)
	if player == nil {
		http.Error(w, "failed to fetch player data", http.StatusNotFound)
		return nil, fmt.Errorf("player data unavailable")
	}

	// 3. Upgrade connection
	conn, err := f.upgrader.Upgrade(w, r, nil)
	if err != nil {
		f.logger.Error("failed to upgrade connection", "error", err)
		return nil, err
	}
	wsConn := connection.WrapConnection(conn)

	// 4. Check for existing match
	address := f.getMatchAddress(ctx, *playerID)
	if address != nil {
		resp := &mmeJson.Response{
			Type:    mmeJson.Match,
			Address: address,
		}
		if err := wsConn.Write(resp); err != nil {
			f.logger.Error("failed to send match response", "error", err)
		}
		_ = wsConn.Close()
		return nil, nil
	}

	// 5. Create client
	clientCtx, clientCancel := context.WithCancel(ctx)
	c := &clientImpl{
		conn:                    wsConn,
		disconnectCh:            f.disconnectCh,
		ctx:                     clientCtx,
		cancel:                  clientCancel,
		checkPoolPresenceTicker: time.NewTicker(f.config.CheckInPoolPeriod),
		checkMatchTicker:        time.NewTicker(f.config.CheckMatchPeriod),
		updateConnectionTicket:  time.NewTicker(f.config.ConnectionTTL),
		state:                   Stale,
		mmPoolRepo:              f.mmPoolRepo,
		logger:                  f.logger,
		player:                  player,
		msClient:                f.msClient,
		config:                  f.config,
	}

	return c, nil
}
