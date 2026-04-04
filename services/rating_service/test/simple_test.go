package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "team_dynamics/api/proto/rating_service"
	"testing"
)

func TestSimple(t *testing.T) {
	conn, err := grpc.NewClient("localhost:50001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal("Cannot create connection")
	}
	defer func() {
		_ = conn.Close()
	}()
	client := pb.NewRatingServiceClient(conn)
	var id1 int64 = 1
	var id2 int64 = 2
	matchResult := pb.MatchResult_MATCH_RESULT_WINNER
	matchId := "hgiohrriuhkrjbhijen"
	request1 := &pb.GetRatingRequest{
		UserId: &id1,
	}
	response1, err := client.GetRating(context.TODO(), request1)
	if err != nil {
		t.Fatalf("Request 1 failed: %v", err)
	}
	t.Logf("Request 1 response: %v\n", response1)

	request2 := &pb.UpdateRatingRequest{
		Player1Id:   &id1,
		Player2Id:   &id2,
		MatchResult: matchResult,
		MatchId:     &matchId,
	}
	response2, err := client.UpdateRating(context.TODO(), request2)
	if err != nil {
		t.Fatalf("Request 2 failed: %v", err)
	}
	t.Logf("Request 2 succeded: %v", response2)

	request3 := &pb.GetRatingRequest{
		UserId: &id2,
	}
	response3, err := client.GetRating(context.TODO(), request3)
	if err != nil {
		t.Fatalf("Request 3 failed: %v", err)
	}
	t.Logf("Request 3 response: %v\n", response3)

	request4 := &pb.UpdateRatingRequest{
		Player1Id:   &id1,
		Player2Id:   &id2,
		MatchResult: matchResult,
		MatchId:     &matchId,
	}
	response4, err := client.UpdateRating(context.TODO(), request4)
	if err != nil {
		t.Fatalf("Request 4 failed: %v", err)
	}
	t.Logf("Request 4 response: %v", response4)
}
