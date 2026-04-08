package services

import (
	allocv1 "agones.dev/agones/pkg/apis/allocation/v1"
	"agones.dev/agones/pkg/client/clientset/versioned"
	"context"
	"errors"
	"log"
	"strconv"
	pb "team_dynamics/api/proto/fleet_manager"
	"team_dynamics/fleet_manager/config"
	"team_dynamics/fleet_manager/k8s"
)

type FleetManagerServiceDep struct {
	K8sClient *versioned.Clientset
	FmConfig  *config.FleetManagerConfig
}

func (s *FleetManagerServiceDep) Allocate(ctx context.Context, request *pb.AllocateRequest) (*pb.AllocateResponse, error) {
	if request.Player1 == nil || request.Player2 == nil {
		return nil, errors.New("one of players is nil")
	}
	for {
		res, err := k8s.Allocate(
			s.K8sClient,
			s.FmConfig,
			map[string]string{
				"player1": strconv.FormatUint(request.GetPlayer1(), 10),
				"player2": strconv.FormatUint(request.GetPlayer2(), 10)},
			ctx,
		)
		if err != nil {
			log.Printf("K8s error: %s", err.Error())
			return nil, err
		}
		switch res.State {
		case allocv1.GameServerAllocationContention:
			log.Printf("Contention happened: retrying")
			continue
		case allocv1.GameServerAllocationUnAllocated:
			log.Printf("Game servers are full")
			return nil, errors.New("game servers are full")
		}
		address := k8s.GetAddress(res.Address, res.Ports)
		if address == nil {
			log.Printf("Unable to allocate game port")
			return nil, errors.New("unable to allocate game port")
		}
		name := res.GameServerName
		return &pb.AllocateResponse{Address: address, Name: &name}, nil
	}
}

func (s *FleetManagerServiceDep) GetServer(ctx context.Context, request *pb.GetServerRequest) (*pb.GetServerResponse, error) {
	res, err := k8s.GetServer(s.K8sClient, s.FmConfig, request.GetName(), ctx)
	if err != nil {
		log.Printf("K8s error: %s", err.Error())
		return nil, err
	}
	address := k8s.GetAddress(res.Status.Address, res.Status.Ports)
	if address == nil {
		log.Print("Unable to find game port")
		return nil, errors.New("unable to find game port")
	}
	return &pb.GetServerResponse{Address: address}, nil
}
