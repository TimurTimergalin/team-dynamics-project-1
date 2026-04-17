package test

import (
	"context"
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"net"
	fmPb "team_dynamics/api/proto/fleet_manager"
	mhsPb "team_dynamics/api/proto/match_history_service"
	pb "team_dynamics/api/proto/match_service"
	"team_dynamics/match_service/controllers"
	msRedis "team_dynamics/match_service/redis"
	"team_dynamics/match_service/services"
	"testing"
	"time"

	rsPb "team_dynamics/api/proto/rating_service"
)

type SimpleTestRsMockClient struct {
	updateRatingCalled int
}

func (*SimpleTestRsMockClient) GetRating(ctx context.Context, in *rsPb.GetRatingRequest, opts ...grpc.CallOption) (*rsPb.GetRatingResponse, error) {
	panic("implement me")
}
func (s *SimpleTestRsMockClient) UpdateRating(ctx context.Context, in *rsPb.UpdateRatingRequest, opts ...grpc.CallOption) (*rsPb.UpdateRatingResponse, error) {
	s.updateRatingCalled += 1
	return &rsPb.UpdateRatingResponse{}, nil
}

type SimpleTestMhsMockClient struct {
	saveMatchCalled int
}

func (s *SimpleTestMhsMockClient) GetMatchHistory(ctx context.Context, in *mhsPb.GetMatchHistoryRequest, opts ...grpc.CallOption) (*mhsPb.GetMatchHistoryResponse, error) {
	panic("implement me")
}
func (s *SimpleTestMhsMockClient) SaveMatch(ctx context.Context, in *mhsPb.SaveMatchRequest, opts ...grpc.CallOption) (*mhsPb.SaveMatchResponse, error) {
	s.saveMatchCalled += 1
	return &mhsPb.SaveMatchResponse{}, nil
}

type SimpleTestFleetManagerMockClient struct {
	allocateCalled  int
	getServerCalled int
}

func (s *SimpleTestFleetManagerMockClient) Allocate(ctx context.Context, in *fmPb.AllocateRequest, opts ...grpc.CallOption) (*fmPb.AllocateResponse, error) {
	s.allocateCalled += 1
	return &fmPb.AllocateResponse{
		ConnectionInfo: &fmPb.ConnectionInfo{
			Address: in.MatchId,
		},
	}, nil
}

func (s *SimpleTestFleetManagerMockClient) GetServer(ctx context.Context, in *fmPb.GetServerRequest, opts ...grpc.CallOption) (*fmPb.GetServerResponse, error) {
	s.getServerCalled += 1
	return &fmPb.GetServerResponse{
		ConnectionInfo: &fmPb.ConnectionInfo{
			Address: in.MatchId,
		},
	}, nil
}

func ptr[T any](v T) *T {
	return &v
}

func ptr64(v int64) *int64 {
	return &v
}

func TestSimple(t *testing.T) {
	mr := miniredis.RunT(t)
	defer mr.Close()

	redisOptions := &redis.Options{
		Addr:         mr.Addr(),
		PoolSize:     1,
		MinIdleConns: 1,
		DialTimeout:  1 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		PoolTimeout:  1 * time.Second,
	}

	rdb := redis.NewClient(redisOptions)
	defer func(rdb *redis.Client) {
		_ = rdb.Close()
	}(rdb)

	rsClient := &SimpleTestRsMockClient{}
	fmClient := &SimpleTestFleetManagerMockClient{}
	mhsClient := &SimpleTestMhsMockClient{}

	controller := &controllers.MatchServiceController{
		Service: services.MakeMatchService(
			msRedis.MakeMatchKvRepo(rdb),
			fmClient,
			rsClient,
			mhsClient,
		),
	}

	lis := bufconn.Listen(1024 * 1024)
	s := grpc.NewServer()
	pb.RegisterMatchServiceServer(s, controller)
	go func() {
		if err := s.Serve(lis); err != nil {
			t.Errorf("server exited: %v", err)
		}
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}
	conn, err := grpc.NewClient("passthrough://bufnet",
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = conn.Close()
		s.Stop()
		_ = lis.Close()
	}()

	testClient := pb.NewMatchServiceClient(conn)

	player1Data1 := &pb.PlayerData{
		PlayerId:     ptr64(1),
		PlayerName:   ptr("Player 1"),
		PlayerRating: ptr64(1200),
		RegId:        ptr("player 1 reg_id 1"),
	}
	player2Data1 := &pb.PlayerData{
		PlayerId:     ptr64(2),
		PlayerName:   ptr("Player 2"),
		PlayerRating: ptr64(1300),
		RegId:        ptr("player 2 reg_id 1"),
	}
	player1Data2 := &pb.PlayerData{
		PlayerId:     ptr64(1),
		PlayerName:   ptr("Player 1"),
		PlayerRating: ptr64(1200),
		RegId:        ptr("player 1 reg_id 2"),
	}
	player2Data2 := &pb.PlayerData{
		PlayerId:     ptr64(2),
		PlayerName:   ptr("Player 2"),
		PlayerRating: ptr64(1300),
		RegId:        ptr("player 2 reg_id 2"),
	}
	player3Data := &pb.PlayerData{
		PlayerId:     ptr64(3),
		PlayerName:   ptr("Player 3"),
		PlayerRating: ptr64(900),
		RegId:        ptr("player 3 reg_id"),
	}
	player1Data3 := &pb.PlayerData{
		PlayerId:     ptr64(1),
		PlayerName:   ptr("Player 1"),
		PlayerRating: ptr64(1200),
		RegId:        ptr("player 1 reg_id 3"),
	}
	player2Data3 := &pb.PlayerData{
		PlayerId:     ptr64(2),
		PlayerName:   ptr("Player 2"),
		PlayerRating: ptr64(1300),
		RegId:        ptr("player 2 reg_id 3"),
	}

	inputMatch1 := &pb.InputMatch{
		Player1: player1Data1,
		Player2: player2Data1,
		Fleet:   ptr("fleet"),
	}
	inputMatch2 := &pb.InputMatch{
		Player1: player1Data2,
		Player2: player2Data2,
		Fleet:   ptr("fleet"),
	}
	inputMatch3 := &pb.InputMatch{
		Player1: player1Data1,
		Player2: player3Data,
		Fleet:   ptr("fleet"),
	}
	inputMatch4 := &pb.InputMatch{
		Player1: player1Data3,
		Player2: player2Data3,
		Fleet:   ptr("fleet"),
	}

	crRequest1 := &pb.StartMatchRequest{
		Matches: []*pb.InputMatch{
			inputMatch1,
		},
	}
	crRequest2 := &pb.StartMatchRequest{
		Matches: []*pb.InputMatch{
			inputMatch1,
			inputMatch3,
		},
	}
	crRequest3 := &pb.StartMatchRequest{
		Matches: []*pb.InputMatch{
			inputMatch2,
		},
	}
	crRequest4 := &pb.StartMatchRequest{
		Matches: []*pb.InputMatch{
			inputMatch4,
		},
	}

	ctx := context.Background()
	response1, err1 := testClient.StartMatch(ctx, crRequest1)
	if err1 != nil {
		t.Log(mr.Dump())
		t.Fatalf("Request 1 failed: %v", err1)
	}
	if response1.Results[0].MatchId == nil {
		t.Log(mr.Dump())
		t.Fatalf("Request 1 response empty: %v", response1)
	}
	matchId1 := response1.Results[0].GetMatchId()
	response2, err2 := testClient.StartMatch(ctx, crRequest2)
	if err2 != nil {
		t.Log(mr.Dump())
		t.Fatalf("Request 2 failed: %v", err2)
	}
	correct := response2.Results[0].Player1FailResponse == pb.PlayerFailResponse_PLAYER_FAIL_RESPONSE_REMOVE &&
		response2.Results[0].Player2FailResponse == pb.PlayerFailResponse_PLAYER_FAIL_RESPONSE_REMOVE &&
		response2.Results[1].Player1FailResponse == pb.PlayerFailResponse_PLAYER_FAIL_RESPONSE_REMOVE &&
		response2.Results[1].Player2FailResponse == pb.PlayerFailResponse_PLAYER_FAIL_RESPONSE_REENTER
	if !correct {
		t.Log(mr.Dump())
		t.Fatalf("Request 2 is incorrect: %v", response2)
	}

	getReq1 := &pb.GetMatchRequest{
		PlayerId: ptr64(3),
	}
	response3, err3 := testClient.GetMatch(ctx, getReq1)
	if err3 != nil {
		t.Log(mr.Dump())
		t.Fatalf("Request 3 failed: %v", err3)
	}
	correct = response3.ConnectionInfo == nil
	if !correct {
		t.Log(mr.Dump())
		t.Fatalf("Request 3 is incorrect: %v", response3)
	}
	cancReq1 := &pb.CancelMatchRequest{
		MatchId: &matchId1,
	}
	_, err4 := testClient.CancelMatch(ctx, cancReq1)
	if err4 != nil {
		t.Log(mr.Dump())
		t.Fatalf("Request 4 failed: %v", err4)
	}
	response5, err5 := testClient.StartMatch(ctx, crRequest3)
	if err5 != nil {
		t.Log(mr.Dump())
		t.Fatalf("Request 5 failed: %v", err5)
	}

	correct = response5.Results[0].MatchId != nil
	if !correct {
		t.Log(mr.Dump())
		t.Fatalf("Request 5 is incorrect: %v", response5)
	}
	matchId2 := response5.Results[0].GetMatchId()

	getReq2 := &pb.GetMatchRequest{
		PlayerId: ptr64(1),
	}

	response6, err6 := testClient.GetMatch(ctx, getReq2)
	if err6 != nil {
		t.Log(mr.Dump())
		t.Fatalf("Request 6 failed: %v", err6)
	}
	correct = response6.ConnectionInfo != nil
	if !correct {
		t.Log(mr.Dump())
		t.Fatalf("Request 6 is incorrect: %v", response6)
	}
	address := response6.ConnectionInfo.GetAddress()
	response7, err7 := testClient.GetMatch(ctx, getReq2)
	if err7 != nil {
		t.Log(mr.Dump())
		t.Fatalf("Request 7 failed: %v", err7)
	}
	correct = response7.ConnectionInfo != nil && response7.ConnectionInfo.GetAddress() == address
	if !correct {
		t.Log(mr.Dump())
		t.Fatalf("Request 7 is incorrect: %v", response7)
	}
	cancReq2 := &pb.CancelMatchRequest{
		MatchId: &matchId2,
	}
	_, err8 := testClient.CancelMatch(ctx, cancReq2)
	if err8 != nil {
		t.Log(mr.Dump())
		t.Fatalf("Request 8 failed: %v", err8)
	}
	response9, err9 := testClient.StartMatch(ctx, crRequest3)
	if err9 != nil {
		t.Log(mr.Dump())
		t.Fatalf("Request 9 failed: %v", err9)
	}
	correct = response9.Results[0].MatchId != nil
	if !correct {
		t.Log(mr.Dump())
		t.Fatalf("Request 9 is incorrect: %v", response9)
	}
	matchId3 := response9.Results[0].GetMatchId()
	response10, err10 := testClient.GetMatch(ctx, getReq2)
	if err10 != nil {
		t.Log(mr.Dump())
		t.Fatalf("Request 10 failed: %v", err10)
	}
	correct = response10.ConnectionInfo != nil
	if !correct {
		t.Log(mr.Dump())
		t.Fatalf("Request 10 is incorrect: %v", response10)
	}
	finReq1 := &pb.EndMatchRequest{
		MatchId:  &matchId3,
		WinnerId: ptr64(1),
	}
	_, err11 := testClient.EndMatch(ctx, finReq1)
	if err11 != nil {
		t.Log(mr.Dump())
		t.Fatalf("Request 11 failed: %v", err11)
	}
	response12, err12 := testClient.StartMatch(ctx, crRequest4)
	if err12 != nil {
		t.Log(mr.Dump())
		t.Fatalf("Request 12 failed: %v", err12)
	}
	correct = response12.Results[0].MatchId != nil
	if !correct {
		t.Log(mr.Dump())
		t.Fatalf("Request 12 is incorrect: %v", response12)
	}
	matchId4 := response12.Results[0].GetMatchId()
	response13, err13 := testClient.GetMatch(ctx, getReq2)
	if err13 != nil {
		t.Log(mr.Dump())
		t.Fatalf("Request 13 failed: %v", err13)
	}
	correct = response13.ConnectionInfo != nil
	if !correct {
		t.Log(mr.Dump())
		t.Fatalf("Request 13 is incorrect: %v", response13)
	}
	finReq2 := &pb.EndMatchRequest{
		MatchId:  &matchId4,
		WinnerId: ptr64(1),
	}
	_, err14 := testClient.EndMatch(ctx, finReq2)
	if err14 != nil {
		t.Log(mr.Dump())
		t.Fatalf("Request 14 failed: %v", err14)
	}
	restReq1 := &pb.RenewMatchRequest{
		MatchId: &matchId4,
	}
	_, err15 := testClient.RenewMatch(ctx, restReq1)
	if err15 != nil {
		t.Log(mr.Dump())
		t.Fatalf("Request 15 failed: %v", err15)
	}
	response16, err16 := testClient.GetMatch(ctx, getReq2)
	if err16 != nil {
		t.Log(mr.Dump())
		t.Fatalf("Request 16 failed: %v", err16)
	}
	correct = response16.ConnectionInfo != nil
	if !correct {
		t.Log(mr.Dump())
		t.Fatalf("Request 16 is incorrect: %v", response16)
	}

	t.Logf("%v", rsClient)
	t.Logf("%v", fmClient)
	t.Logf("%v", mhsClient)
}
