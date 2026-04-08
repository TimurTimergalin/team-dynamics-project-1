package services

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "team_dynamics/api/proto/fleet_manager"
	"team_dynamics/fleet_manager/k8s"
)

type FleetManagerService interface {
	Allocate(ctx context.Context, request *pb.AllocateRequest) (*pb.AllocateResponse, error)
	GetServer(ctx context.Context, request *pb.GetServerRequest) (*pb.GetServerResponse, error)
}

type fleetManagerServiceImpl struct {
	k8sOps k8s.Ops
}

func MakeFleetManagerService(k8sOps k8s.Ops) FleetManagerService {
	return &fleetManagerServiceImpl{k8sOps}
}

func validatePlayerAnnotations(player *pb.PlayerAnnotations) error {
	if player == nil {
		return errors.New("no player")
	}
	if player.Id == nil {
		return errors.New("no player id")
	}
	if player.Name == nil {
		return errors.New("no player name")
	}
	if player.Rating == nil {
		return errors.New("no player rating")
	}
	return nil
}

func validateAllocateRequest(request *pb.AllocateRequest) error {
	if err := validatePlayerAnnotations(request.Player1); err != nil {
		return fmt.Errorf("while validating player1: %w", err)
	}
	if err := validatePlayerAnnotations(request.Player2); err != nil {
		return fmt.Errorf("while validating player2: %w", err)
	}
	if request.MatchId == nil {
		return errors.New("no match id")
	}
	return nil
}

func validateGetServerRequest(request *pb.GetServerRequest) error {
	if request.MatchId == nil {
		return fmt.Errorf("no match id")
	}
	return nil
}

func makeAnnotations(matchId string, player1, player2 *pb.PlayerAnnotations) map[string]string {
	return map[string]string{
		"player1_id":     player1.GetId(),
		"player1_name":   player1.GetName(),
		"player1_rating": player1.GetRating(),
		"player2_id":     player2.GetId(),
		"player2_name":   player2.GetName(),
		"player2_rating": player2.GetRating(),
		"match_id":       matchId,
	}
}

func convertError(opsErr *k8s.OpsError) error {
	if opsErr == nil {
		return nil
	}
	switch opsErr.Type {
	case k8s.NotFoundError:
		return status.Error(codes.NotFound, opsErr.Error())
	case k8s.ContentionError:
	case k8s.ConnectionError:
	case k8s.FleetFullError:
		return status.Error(codes.Unavailable, opsErr.Error())
	default:
	}
	return status.Error(codes.Unknown, opsErr.Error())
}

func makeConnectionInfo(address string) *pb.ConnectionInfo {
	return &pb.ConnectionInfo{
		Address: &address,
	}
}

func (s *fleetManagerServiceImpl) Allocate(ctx context.Context, request *pb.AllocateRequest) (*pb.AllocateResponse, error) {
	if err := validateAllocateRequest(request); err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid request: %v", err))
	}
	address, err := s.k8sOps.GetServerByMatchId(ctx, request.GetMatchId())
	if err == nil {
		return &pb.AllocateResponse{
			ConnectionInfo: makeConnectionInfo(address),
		}, nil
	}
	annotations := makeAnnotations(request.GetMatchId(), request.Player1, request.GetPlayer2())
	address, err = s.k8sOps.Allocate(ctx, request.GetMatchId(), request.FleetName, annotations)
	if err == nil {
		return &pb.AllocateResponse{
			ConnectionInfo: makeConnectionInfo(address),
		}, nil
	}
	switch err.Type {
	case k8s.FleetFullError:
		if request.FleetName != nil {
			address, err = s.k8sOps.Allocate(ctx, request.GetMatchId(), nil, annotations)
			if err == nil {
				return &pb.AllocateResponse{
					ConnectionInfo: makeConnectionInfo(address),
				}, nil
			}
		}
	default:
	}
	return nil, convertError(err)
}

func (s *fleetManagerServiceImpl) GetServer(ctx context.Context, request *pb.GetServerRequest) (*pb.GetServerResponse, error) {
	if err := validateGetServerRequest(request); err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid request: %v", err))
	}
	address, err := s.k8sOps.GetServerByMatchId(ctx, request.GetMatchId())
	if err != nil {
		return nil, convertError(err)
	}
	return &pb.GetServerResponse{
		ConnectionInfo: makeConnectionInfo(address),
	}, nil
}
