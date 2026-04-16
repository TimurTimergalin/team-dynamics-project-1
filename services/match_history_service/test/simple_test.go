package test

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"strconv"
	pb "team_dynamics/api/proto/match_history_service"
	"testing"
	"time"
)

func TestSimple(t *testing.T) {
	conn, err := grpc.NewClient("localhost:50001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal("Cannot create connection")
	}
	defer func() {
		_ = conn.Close()
	}()

	client := pb.NewMatchHistoryServiceClient(conn)
	var id1 int64 = 1
	var id2 int64 = 2
	name1 := "name1"
	name2 := "name2"
	var rating1 int64 = 1200
	var rating2 int64 = 1100
	player1 := &pb.ParticipantData{
		Id:     &id1,
		Name:   &name1,
		Rating: &rating1,
	}
	player2 := &pb.ParticipantData{
		Id:     &id2,
		Name:   &name2,
		Rating: &rating2,
	}

	baseTimestamp := time.Now().UTC().UnixMilli()
	var lastRequest *pb.SaveMatchRequest

	for i := int64(0); i < 18; i += 1 {
		var player1First = i%2 == 0
		var neg = !player1First
		time1 := 15*10000 + i*5*1000
		time2 := time1 + 3000
		endTime := baseTimestamp + 100000*i
		var result pb.MatchResult
		if i%2 == 0 {
			result = pb.MatchResult_MATCH_RESULT_PLAYER2_WIN
		} else {
			result = pb.MatchResult_MATCH_RESULT_PLAYER1_WIN
		}
		if i%3 == 0 {
			result = pb.MatchResult_MATCH_RESULT_CANCELLED
		}
		matchid := strconv.FormatInt(i, 10)
		request := &pb.SaveMatchRequest{
			Match: &pb.MatchData{
				MatchId: &matchid,
				Player1: player1,
				Player2: player2,
				Rounds: []*pb.Round{
					{
						IsPlayer1Killer: &player1First,
						TimeMillis:      &time1,
					},
					{
						IsPlayer1Killer: &neg,
						TimeMillis:      &time2,
					},
				},
				MatchResult:  result,
				EndTimestamp: &endTime,
			},
		}
		lastRequest = request
		_, err := client.SaveMatch(context.Background(), request)
		if err != nil {
			t.Fatalf("Error while saving match %d: %v", i, err)
		}
	}

	response, err := client.GetMatchHistory(context.Background(), &pb.GetMatchHistoryRequest{
		UserId: &id1,
	})
	if err != nil {
		t.Fatalf("first page response failed: %v", err)
	}
	if len(response.Matches) != 10 {
		t.Fatalf("wrong matches count in first page: %d != 10", len(response.Matches))
	}
	if response.Pagekey == nil {
		t.Fatalf("no pagekey in first page response")
	}
	t.Logf("first page response: %v", response)
	response, err = client.GetMatchHistory(context.Background(), &pb.GetMatchHistoryRequest{
		UserId:  &id1,
		Pagekey: response.Pagekey,
	})
	if err != nil {
		t.Fatalf("second page request failed: %v", err)
	}
	if len(response.Matches) != 8 {
		t.Fatalf("wrong matches count in second page: %d != 8", len(response.Matches))
	}
	if response.Pagekey != nil {
		t.Fatalf("page key is present even though it shouldn't be: %v", response.Pagekey)
	}
	t.Logf("second page response: %v", response)

	_, err = client.SaveMatch(context.Background(), lastRequest)
	if err != nil {
		t.Fatalf("Repeated request counts as an error: %v", err)
	}
}
