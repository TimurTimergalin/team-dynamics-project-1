package test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/alicebob/miniredis/v2"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"runtime/debug"
	"sync"
	fmPb "team_dynamics/api/proto/fleet_manager"
	msPb "team_dynamics/api/proto/match_service"
	rsPb "team_dynamics/api/proto/rating_service"
	usPb "team_dynamics/api/proto/user_service"
	"team_dynamics/user_event/client"
	"team_dynamics/user_event/controllers"
	"team_dynamics/user_event/downstream"
	ueJson "team_dynamics/user_event/json"
	ueRedis "team_dynamics/user_event/redis"
	"testing"
	"time"
)

func ptr[T any](v T) *T {
	return &v
}

type mockMatchServiceFactory struct{}

func (m *mockMatchServiceFactory) GetMatch(_ context.Context, _ *msPb.GetMatchRequest) (*msPb.GetMatchResponse, error) {
	return &msPb.GetMatchResponse{ConnectionInfo: &fmPb.ConnectionInfo{Address: ptr("some_address")}, Player1Id: ptr[int64](1), Player2Id: ptr[int64](2)}, nil
}

func (m *mockMatchServiceFactory) StartMatch(_ context.Context, _ *msPb.StartMatchRequest) (*msPb.StartMatchResponse, error) {
	return &msPb.StartMatchResponse{Results: []*msPb.MatchCreationResult{
		{MatchId: ptr("1")},
	}}, nil
}

type mockRatingServiceFactory struct{}

func (m *mockRatingServiceFactory) GetRating(_ context.Context, _ *rsPb.GetRatingRequest) (*rsPb.GetRatingResponse, error) {
	return &rsPb.GetRatingResponse{Rating: &rsPb.RatingData{RatingValue: ptr(1200.), DisplayValue: ptr[int64](1200)}}, nil
}

type mockUserServiceFactory struct{}

func (m *mockUserServiceFactory) GetUserData(_ context.Context, req *usPb.GetUserDataRequest) (*usPb.GetUserDataResponse, error) {
	return &usPb.GetUserDataResponse{UserData: &usPb.UserData{Id: req.Id, Name: ptr(fmt.Sprintf("Player %d", *req.Id))}}, nil
}

var _ downstream.MatchServiceClientFactory = (*mockMatchServiceFactory)(nil)
var _ downstream.RatingServiceClientFactory = (*mockRatingServiceFactory)(nil)
var _ downstream.UserServiceClientFactory = (*mockUserServiceFactory)(nil)

type TestClient struct {
	conn *websocket.Conn
	mu   sync.Mutex
}

func NewTestClient(serverAddr, playerID string) (*TestClient, error) {
	u := url.URL{Scheme: "ws", Host: serverAddr, Path: "/events"}
	q := u.Query()
	q.Set("playerId", playerID)
	u.RawQuery = q.Encode()
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}
	return &TestClient{conn: conn}, nil
}

func (c *TestClient) SendRequest(req *ueJson.Request) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.conn.WriteJSON(req)
}

func (c *TestClient) ReadResponse() (*ueJson.Response, error) {
	var resp ueJson.Response
	_ = c.conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	defer func(conn *websocket.Conn, t time.Time) {
		_ = conn.SetReadDeadline(t)
	}(c.conn, *new(time.Time))
	err := c.conn.ReadJSON(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *TestClient) Close() error {
	return c.conn.Close()
}

func TestSimple(t *testing.T) {
	mr := miniredis.RunT(t)
	defer mr.Close()

	redisOptions := &redis.Options{
		Addr:         mr.Addr(),
		PoolSize:     2,
		MinIdleConns: 2,
		DialTimeout:  2 * time.Second,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
		PoolTimeout:  2 * time.Second,
	}

	rdb := redis.NewClient(redisOptions)
	defer func(rdb *redis.Client) {
		_ = rdb.Close()
	}(rdb)

	rsClient := &mockRatingServiceFactory{}
	msClient := &mockMatchServiceFactory{}
	usClient := &mockUserServiceFactory{}

	userKvRepo := ueRedis.MakeUserKvRepo(rdb)
	registerCh := make(chan client.Client, 32)
	disconnectCh := make(chan client.Client, 32)
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug})
	logger := slog.New(handler)
	upgrader := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(*http.Request) bool {
			return true
		},
	}
	clientCfg := &client.UserEventConfig{
		CheckStatusPeriod: 1 * time.Second,
	}
	ctx := context.Background()
	hub := client.NewHub(disconnectCh, registerCh, logger, clientCfg)
	factory := client.NewClientFactory(usClient, rsClient, msClient, userKvRepo, disconnectCh, logger, clientCfg, upgrader)
	controller := controllers.NewUserEventController(ctx, hub, factory, "localhost:8000")

	assert := func(cond bool, msg string, args ...any) {
		if !cond {
			t.Log(mr.Dump())
			debug.PrintStack()
			t.Fatalf(msg, args...)
		}
	}

	go controller.Run(ctx)
	t.Log("Testing connect 1")
	cl1, err := NewTestClient("localhost:8000", "1")
	assert(err == nil, "cannot connect to server 1: %v", err)

	t.Log("Testing repeated connection")
	_, err = NewTestClient("localhost:8000", "1")
	assert(err != nil, "second connection succeeded")

	t.Log("Testing subscription on offline")
	req1 := &ueJson.Request{
		Type: ueJson.Subscribe,
		Payload: ueJson.SubscribePayload{
			Users: []int64{
				2,
			},
		},
	}
	err = cl1.SendRequest(req1)
	assert(err == nil, "Unable to send subscribe 1: %v", err)

	resp, err := cl1.ReadResponse()
	assert(err == nil, "Unable to receive subscribe response 1: %v", err)

	assert(len(resp.Events) == 1, "Unknown events")
	ev := resp.Events[0]
	assert(ev.Type == ueJson.StatusChanged, "wrong event")

	{
		var payload ueJson.StatusChangedPayload
		assert(json.Unmarshal(ev.Payload, &payload) == nil, "Wrong payload")
		assert(payload.UserId == 2, "Wrong user_id")
		assert(payload.NewStatus == ueJson.Offline, "Wrong status")
	}

	t.Log("Testing connection 2")
	cl2, err := NewTestClient("localhost:8000", "2")
	assert(err == nil, "cannot connect to server 2: %v", err)

	t.Log("Test update on login")
	resp, err = cl1.ReadResponse()
	assert(err == nil, "Unable to receive subscribe response 2: %v", err)
	assert(len(resp.Events) == 1, "Unknown events")
	ev = resp.Events[0]
	assert(ev.Type == ueJson.StatusChanged, "wrong event")
	{
		var payload ueJson.StatusChangedPayload
		assert(json.Unmarshal(ev.Payload, &payload) == nil, "Wrong payload")
		assert(payload.UserId == 2, "Wrong user_id")
		assert(payload.NewStatus == ueJson.Online, "Wrong status")
	}

	t.Log("Testing notify busy")
	err = cl1.SendRequest(&ueJson.Request{
		Type:    ueJson.NotifyBusy,
		Payload: ueJson.NotifyBusyPayload{},
	})
	assert(err == nil, "Unable to send request")
	time.Sleep(100 * time.Millisecond)
	err = cl2.SendRequest(&ueJson.Request{
		Type: ueJson.Subscribe,
		Payload: ueJson.SubscribePayload{
			Users: []int64{
				1,
			},
		},
	})
	assert(err == nil, "Unable to send request")

	resp, err = cl2.ReadResponse()
	assert(err == nil, "Unable to read response: %v", err)
	assert(len(resp.Events) == 1, "Unknown events: %v", resp)
	ev = resp.Events[0]
	assert(ev.Type == ueJson.StatusChanged, "Wrong status")
	{
		var payload ueJson.StatusChangedPayload
		assert(json.Unmarshal(ev.Payload, &payload) == nil, "Wrong payload")
		assert(payload.UserId == 1, "Wrong user_id")
		assert(payload.NewStatus == ueJson.InGame, "Wrong status: %v", payload.NewStatus)
	}

	t.Log("Testing notify free")
	err = cl1.SendRequest(&ueJson.Request{
		Type:    ueJson.NotifyFree,
		Payload: ueJson.NotifyFreePayload{},
	})
	assert(err == nil, "Unable to send request")
	resp, err = cl2.ReadResponse()
	assert(err == nil, "Unable to read response: %v", err)
	assert(len(resp.Events) == 1, "Unknown events: %v", resp)
	ev = resp.Events[0]
	assert(ev.Type == ueJson.StatusChanged, "Wrong status")
	{
		var payload ueJson.StatusChangedPayload
		assert(json.Unmarshal(ev.Payload, &payload) == nil, "Wrong payload")
		assert(payload.UserId == 1, "Wrong user_id")
		assert(payload.NewStatus == ueJson.Online, "Wrong status")
	}

	t.Log("Sending challenge for testing cancel")
	err = cl1.SendRequest(&ueJson.Request{
		Type: ueJson.Challenge,
		Payload: ueJson.ChallengePayload{
			UserId: 2,
		},
	})
	assert(err == nil, "Unable to send request")
	resp, err = cl2.ReadResponse()
	assert(err == nil, "Unable to read response: %v", err)
	assert(len(resp.Events) == 1, "Unknown events: %v", resp)
	ev = resp.Events[0]
	assert(ev.Type == ueJson.ChallengeReceived, "Wrong type")
	{
		var payload ueJson.ChallengeReceivedPayload
		assert(json.Unmarshal(ev.Payload, &payload) == nil, "Wrong payload")
		assert(payload.UserId == 1, "Wrong id")
		assert(payload.UserName == "Player 1", "Wrong name")
		assert(payload.MessageId != "", "Match Id absent")
	}

	t.Log("Testing cancelling")
	err = cl1.SendRequest(&ueJson.Request{
		Type:    ueJson.CancelChallenge,
		Payload: ueJson.CancelChallengePayload{},
	})
	resp, err = cl2.ReadResponse()
	assert(err == nil, "Unable to read response: %v", err)
	assert(len(resp.Events) == 1, "Unknown events: %v", resp)
	ev = resp.Events[0]
	assert(ev.Type == ueJson.ChallengeCancelled, "Wrong type")
	{
		var payload ueJson.ChallengeCancelledPayload
		assert(json.Unmarshal(ev.Payload, &payload) == nil, "Wrong payload")
	}

	t.Log("Sending challenge for testing decline")
	err = cl1.SendRequest(&ueJson.Request{
		Type: ueJson.Challenge,
		Payload: ueJson.ChallengePayload{
			UserId: 2,
		},
	})
	assert(err == nil, "Unable to send request")
	resp, err = cl2.ReadResponse()
	assert(err == nil, "Unable to read response: %v", err)
	assert(len(resp.Events) == 1, "Unknown events: %v", resp)
	ev = resp.Events[0]
	assert(ev.Type == ueJson.ChallengeReceived, "Wrong type")
	messageId := ""
	{
		var payload ueJson.ChallengeReceivedPayload
		assert(json.Unmarshal(ev.Payload, &payload) == nil, "Wrong payload")
		assert(payload.UserId == 1, "Wrong id")
		assert(payload.UserName == "Player 1", "Wrong name")
		assert(payload.MessageId != "", "Match Id absent")
		messageId = payload.MessageId
	}

	t.Logf("Testing decline, message id: %v", messageId)
	err = cl2.SendRequest(&ueJson.Request{
		Type: ueJson.DeclineChallenge,
		Payload: ueJson.DeclineChallengePayload{
			MessageId: messageId,
			UserId:    1,
		},
	})
	assert(err == nil, "Unable to send request")
	resp, err = cl1.ReadResponse()
	assert(err == nil, "Unable to read response: %v", err)
	assert(len(resp.Events) == 1, "Unknown events: %v", resp)
	ev = resp.Events[0]
	assert(ev.Type == ueJson.ChallengeDeclined, "Wrong type")
	{
		var payload ueJson.ChallengeDeclinedPayload
		assert(json.Unmarshal(ev.Payload, &payload) == nil, "Wrong payload")
	}

	t.Log("Sending challenge for accept")
	err = cl1.SendRequest(&ueJson.Request{
		Type: ueJson.Challenge,
		Payload: ueJson.ChallengePayload{
			UserId: 2,
		},
	})
	assert(err == nil, "Unable to send request")
	resp, err = cl2.ReadResponse()
	assert(err == nil, "Unable to read response: %v", err)
	assert(len(resp.Events) == 1, "Unknown events: %v", resp)
	ev = resp.Events[0]
	assert(ev.Type == ueJson.ChallengeReceived, "Wrong type")
	{
		var payload ueJson.ChallengeReceivedPayload
		assert(json.Unmarshal(ev.Payload, &payload) == nil, "Wrong payload")
		assert(payload.UserId == 1, "Wrong id")
		assert(payload.UserName == "Player 1", "Wrong name")
		assert(payload.MessageId != "", "Match Id absent")
		messageId = payload.MessageId
	}

	t.Logf("Testing accept, message id: %v", messageId)
	err = cl2.SendRequest(&ueJson.Request{
		Type: ueJson.AcceptChallenge,
		Payload: ueJson.AcceptChallengePayload{
			MessageId: messageId,
			UserId:    1,
		},
	})
	assert(err == nil, "Unable to send request")

	for i := 0; i < 2; i++ {
		resp, err = cl2.ReadResponse()
		assert(err == nil, "Unable to read response: %v", err)
		assert(len(resp.Events) == 1, "Unknown events: %v", resp)
		ev = resp.Events[0]
		assert(ev.Type == ueJson.MatchStarted || ev.Type == ueJson.StatusChanged, "Wrong type: %v", ev.Type)
		if ev.Type == ueJson.StatusChanged {
			var payload ueJson.StatusChangedPayload
			assert(json.Unmarshal(ev.Payload, &payload) == nil, "Wrong payload")
			assert(payload.UserId == 1, "Wrong user_id")
			assert(payload.NewStatus == ueJson.InGame, "Wrong status")
		} else {
			var payload ueJson.MatchStartedPayload
			assert(json.Unmarshal(ev.Payload, &payload) == nil, "Wrong payload")
			assert(payload.Address != "", "Address absent")
		}
	}

	resp, err = cl1.ReadResponse()
	assert(err == nil, "Unable to read response: %v", err)
	assert(len(resp.Events) == 2, "Unknown events: %v", resp)
	ev = resp.Events[0]
	assert(ev.Type == ueJson.StatusChanged, "Wrong type")
	{
		var payload ueJson.StatusChangedPayload
		assert(json.Unmarshal(ev.Payload, &payload) == nil, "Wrong payload")
		assert(payload.UserId == 2, "Wrong user_id")
		assert(payload.NewStatus == ueJson.InGame, "Wrong status")
	}
	ev = resp.Events[1]
	assert(ev.Type == ueJson.ChallengeAccepted, "Wrong type")
	{
		var payload ueJson.ChallengeAcceptedPayload
		assert(json.Unmarshal(ev.Payload, &payload) == nil, "Wrong payload")
		assert(payload.Address != "", "Empty address")
	}
}
