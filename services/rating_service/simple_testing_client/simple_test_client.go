package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	pb "team_dynamics/api/proto/rating_service"
)

func main() {
	conn, err := grpc.NewClient("rating-service.rating.svc.cluster.local:50001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Cannot connect: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("Cannot gracefully close connection")
		}
	}(conn)

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
		log.Fatalf("Request 1 failed: %v", err)
	}
	log.Printf("Request 1 response: %v\n", response1)

	request2 := &pb.UpdateRatingRequest{
		Player1Id:   &id1,
		Player2Id:   &id2,
		MatchResult: matchResult,
		MatchId:     &matchId,
	}
	response2, err := client.UpdateRating(context.TODO(), request2)
	if err != nil {
		log.Printf("Request 2 failed: %v", err)
	} else {
		log.Fatalf("Request 2 succeded even though it shouldn't have: %v", response2)
	}

	request3 := &pb.GetRatingRequest{
		UserId: &id2,
	}
	response3, err := client.GetRating(context.TODO(), request3)
	if err != nil {
		log.Fatalf("Request 3 failed: %v", err)
	}
	log.Printf("Request 3 response: %v\n", response3)

	request4 := &pb.UpdateRatingRequest{
		Player1Id:   &id1,
		Player2Id:   &id2,
		MatchResult: matchResult,
		MatchId:     &matchId,
	}
	response4, err := client.UpdateRating(context.TODO(), request4)
	if err != nil {
		log.Fatalf("Request 4 failed: %v", err)
	}
	log.Printf("Request 4 response: %v\n", response4)
	log.Printf("OMG IT WORKS!!!!!!!!!!!")
}
