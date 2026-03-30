package main

import (
	"context"
	"fleet_manager/config"
	pb "fleet_manager/fleet_manager/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func main() {
	fmConfig, err := config.GetConfig()
	if err != nil {
		log.Fatalf("Error while getting fleet manager config, %s", err.Error())
	}

	conn, err := grpc.NewClient(fmConfig.ListenAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Error while establishing connection, %s", err.Error())
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("Unable to close connection gracefully: %s", err.Error())
		}
	}(conn)

	client := pb.NewFleetManagerClient(conn)
	var player1 uint64 = 1
	var player2 uint64 = 2
	allocRequest := &pb.AllocateRequest{Player1: &player1, Player2: &player2}
	allocResponse, err := client.Allocate(context.TODO(), allocRequest)
	if err != nil {
		log.Fatalf("Allocation error: %s", err.Error())
	}

	log.Print(allocResponse.GetName(), allocResponse.GetAddress())
	getAddressRequest := &pb.GetServerRequest{Name: allocResponse.Name}
	getAddressResponse, err := client.GetServer(context.TODO(), getAddressRequest)
	if err != nil {
		log.Fatalf("Get server error: %s", err.Error())
	}
	log.Print(getAddressResponse.GetAddress())
}
