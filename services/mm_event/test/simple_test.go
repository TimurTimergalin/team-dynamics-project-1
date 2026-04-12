package test

import (
	"context"
	"github.com/alicebob/miniredis/v2"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"sync"
	fmPb "team_dynamics/api/proto/fleet_manager"
	msPb "team_dynamics/api/proto/match_service"
	rsPb "team_dynamics/api/proto/rating_service"
	usPb "team_dynamics/api/proto/user_service"
	"team_dynamics/mm_event/client"
	"team_dynamics/mm_event/controllers"
	"team_dynamics/mm_event/json"
	mmeRedis "team_dynamics/mm_event/redis"
	"testing"
	"time"
)

func ptr[T any](v T) *T {
	return &v
}

type SimpleTestMatchServiceMockClient struct {
	hasMatch bool
	mu       sync.Mutex
}

func (s *SimpleTestMatchServiceMockClient) StartMatch(ctx context.Context, in *msPb.StartMatchRequest, opts ...grpc.CallOption) (*msPb.StartMatchResponse, error) {
	panic("implement me")
}

func (s *SimpleTestMatchServiceMockClient) GetMatch(ctx context.Context, in *msPb.GetMatchRequest, opts ...grpc.CallOption) (*msPb.GetMatchResponse, error) {
	if !s.hasMatch {
		return &msPb.GetMatchResponse{}, nil
	}
	return &msPb.GetMatchResponse{ConnectionInfo: &fmPb.ConnectionInfo{Address: ptr("localhost:7777")}}, nil
}

func (s *SimpleTestMatchServiceMockClient) EndMatch(ctx context.Context, in *msPb.EndMatchRequest, opts ...grpc.CallOption) (*msPb.EndMatchResponse, error) {
	panic("implement me")
}

func (s *SimpleTestMatchServiceMockClient) CancelMatch(ctx context.Context, in *msPb.CancelMatchRequest, opts ...grpc.CallOption) (*msPb.CancelMatchResponse, error) {
	panic("implement me")
}

func (s *SimpleTestMatchServiceMockClient) RenewMatch(ctx context.Context, in *msPb.RenewMatchRequest, opts ...grpc.CallOption) (*msPb.RenewMatchResponse, error) {
	panic("implement me")
}

func (s *SimpleTestMatchServiceMockClient) setEnable(v bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.hasMatch = v
}

type SimpleTestRatingServiceMockClient struct{}

func (s *SimpleTestRatingServiceMockClient) GetRating(ctx context.Context, in *rsPb.GetRatingRequest, opts ...grpc.CallOption) (*rsPb.GetRatingResponse, error) {
	return &rsPb.GetRatingResponse{Rating: &rsPb.RatingData{RatingValue: ptr(1200.5), DisplayValue: ptr[int64](1200)}}, nil
}

func (s *SimpleTestRatingServiceMockClient) UpdateRating(ctx context.Context, in *rsPb.UpdateRatingRequest, opts ...grpc.CallOption) (*rsPb.UpdateRatingResponse, error) {
	panic("implement me")
}

type SimpleTestUserServiceMockClient struct{}

func (s *SimpleTestUserServiceMockClient) GetSelfData(ctx context.Context, in *usPb.GetSelfDataRequest, opts ...grpc.CallOption) (*usPb.GetSelfDataResponse, error) {
	panic("implement me")
}

func (s *SimpleTestUserServiceMockClient) GetUserData(ctx context.Context, in *usPb.GetUserDataRequest, opts ...grpc.CallOption) (*usPb.GetUserDataResponse, error) {
	return &usPb.GetUserDataResponse{UserData: &usPb.UserData{Id: ptr[int64](1), Name: ptr("Player1")}}, nil
}

func (s *SimpleTestUserServiceMockClient) GetFriends(ctx context.Context, in *usPb.GetFriendsRequest, opts ...grpc.CallOption) (*usPb.GetFriendsResponse, error) {
	panic("implement me")
}

func (s *SimpleTestUserServiceMockClient) GetIncomingRequests(ctx context.Context, in *usPb.GetIncomingRequestsRequest, opts ...grpc.CallOption) (*usPb.GetIncomingRequestsResponse, error) {
	panic("implement me")
}

func (s *SimpleTestUserServiceMockClient) GetOutgoingRequests(ctx context.Context, in *usPb.GetOutgoingRequestsRequest, opts ...grpc.CallOption) (*usPb.GetOutgoingRequestsResponse, error) {
	panic("implement me")
}

type TestClient struct {
	conn *websocket.Conn
	url  string
	mu   sync.Mutex
}

func NewTestClient(serverAddr, playerID, fleet string) (*TestClient, error) {
	u := url.URL{Scheme: "ws", Host: serverAddr, Path: "/events"}
	q := u.Query()
	q.Set("playerId", playerID)
	q.Set("fleet", fleet)
	u.RawQuery = q.Encode()

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}
	return &TestClient{conn: conn, url: u.String()}, nil
}

func (c *TestClient) SendRequest(req *json.Request) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.conn.WriteJSON(req)
}

func (c *TestClient) ReadResponse() (*json.Response, error) {
	var resp json.Response
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

	rsClient := &SimpleTestRatingServiceMockClient{}
	msClient := &SimpleTestMatchServiceMockClient{}
	usClient := &SimpleTestUserServiceMockClient{}
	mmPoolRepo := mmeRedis.MakeMMPoolRepo(rdb, &mmeRedis.MMPoolRepoConfig{LockTTL: 10 * time.Second})
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

	clientConfig := &client.Config{
		MessageReceivedTimeout: 10 * time.Second,
		CheckMatchPeriod:       1 * time.Second,
		CheckInPoolPeriod:      1 * time.Second,
		HubRegisterTimeout:     10 * time.Second,
	}

	controller := controllers.NewMMEventController(
		client.NewHub(disconnectCh, registerCh, logger, clientConfig),
		client.NewClientFactory(msClient, rsClient, usClient, mmPoolRepo, disconnectCh, logger, clientConfig, upgrader),
		"localhost:8000",
	)

	go controller.Run()
	cl, err := NewTestClient("localhost:8000", "1", "fleet")
	if err != nil {
		t.Log(mr.Dump())
		t.Fatalf("cannot connect to server: %v", err)
	}
	_, err = NewTestClient("localhost:8000", "1", "fleet")
	if err == nil {
		t.Log(mr.Dump())
		t.Fatalf("second connection succeeded")
	}
	t.Logf("reason why second connection failed: %v", err)

	err = cl.SendRequest(&json.Request{Type: json.Register})
	if err != nil {
		t.Log(mr.Dump())
		t.Fatalf("couldn't send register request: %v", err)
	}
	resp, err := cl.ReadResponse()
	if err != nil {
		t.Fatalf("couldn't read register success response: %v", err)
	}
	if resp.Type != json.Registered {
		t.Fatalf("wrong registered response: %v", resp)
	}
	t.Logf("received registered response: %v", resp)
	poolRes, err := rdb.LRange(context.Background(), "mmpool", 0, -1).Result()
	if err != nil {
		t.Log(mr.Dump())
		t.Fatalf("couldn't read pool from redis: %v", err)
	}
	if len(poolRes) != 1 || poolRes[0] != "1" {
		t.Log(mr.Dump())
		t.Fatalf("not in pool")
	}
	keysRes, err := rdb.MGet(context.Background(), "rating:1", "fleet:1", "name:1", "displayed_rating:1", "reg_id:1").Result()
	if err != nil {
		t.Log(mr.Dump())
		t.Fatalf("couldn't read keys from redis: %v", err)
	}
	t.Logf("Res: %v, %v, %v, %v, %v", keysRes...)
	ok := keysRes[0] == "1200.5" && keysRes[1] == "fleet" && keysRes[2] == "Player1" && keysRes[3] == "1200" && keysRes[4] != nil
	if !ok {
		t.Log(mr.Dump())
		t.Fatalf("Wrong keys")
	}
	time.Sleep(4 * time.Second)
	msClient.setEnable(true)
	_, err = rdb.Del(context.Background(), "mmpool", "rating:1", "fleet:1", "name:1", "displayed_rating:1", "reg_id:1").Result()
	if err != nil {
		t.Fatalf("Could not delete keys from redis: %v", err)
	}
	resp, err = cl.ReadResponse()
	if err != nil {
		t.Fatalf("couldn't read address response")
	}
	if resp.Type != json.Match || *resp.Address != "localhost:7777" {
		t.Fatalf("invalid match response: %v", resp)
	}
	resp, err = cl.ReadResponse()
	if err == nil {
		t.Fatalf("Some other response: %v", resp)
	}
	_ = cl.Close()
	cl, err = NewTestClient("localhost:8000", "1", "fleet")
	if err != nil {
		t.Fatalf("could not connect second time: %v", err)
	}
	resp, err = cl.ReadResponse()
	if err != nil {
		t.Fatalf("could not read immediate match response: %v", resp)
	}
	if resp.Type != json.Match || *resp.Address != "localhost:7777" {
		t.Fatalf("invalid immediate match response: %v", resp)
	}
	resp, err = cl.ReadResponse()
	if err == nil {
		t.Fatalf("Some other response: %v", resp)
	}
	_ = cl.Close()
}
