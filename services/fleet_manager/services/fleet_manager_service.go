package services

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "team_dynamics/api/proto/fleet_manager"
	"team_dynamics/fleet_manager/k8s"
	"team_dynamics/logging"
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
		"agones.dev/sdk-player1-id":     player1.GetId(),
		"agones.dev/sdk-player1-name":   player1.GetName(),
		"agones.dev/sdk-player1-rating": player1.GetRating(),
		"agones.dev/sdk-player2-id":     player2.GetId(),
		"agones.dev/sdk-player2-name":   player2.GetName(),
		"agones.dev/sdk-player2-rating": player2.GetRating(),
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
	logger := logging.GetLogger(ctx)
	if err := validateAllocateRequest(request); err != nil {
		logger.Debug("Allocate: invalid request", "error", err)
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid request: %v", err))
	}
	address, _, err := s.k8sOps.GetServerByMatchId(ctx, request.GetMatchId())
	if err == nil {
		logger.Debug("Allocate: server already exists", "matchId", request.GetMatchId(), "address", address)
		return &pb.AllocateResponse{
			ConnectionInfo: makeConnectionInfo(address),
		}, nil
	}
	annotations := makeAnnotations(request.GetMatchId(), request.Player1, request.GetPlayer2())
	address, err = s.k8sOps.Allocate(ctx, request.GetMatchId(), request.FleetName, annotations)
	if err == nil {
		logger.Debug("Allocate: allocated successfully", "matchId", request.GetMatchId(), "address", address)
		return &pb.AllocateResponse{
			ConnectionInfo: makeConnectionInfo(address),
		}, nil
	}
	logger.Debug("Allocate: initial allocation failed", "matchId", request.GetMatchId(), "errorType", err.Type, "error", err)

	tryGettingAgain := false
	switch err.Type {
	case k8s.FleetFullError:
		if request.FleetName != nil {
			logger.Debug("Allocate: fleet full, retrying without fleet constraint", "matchId", request.GetMatchId())
			address, err = s.k8sOps.Allocate(ctx, request.GetMatchId(), nil, annotations)
			if err == nil {
				logger.Debug("Allocate: allocated successfully without fleet constraint", "matchId", request.GetMatchId(), "address", address)
				return &pb.AllocateResponse{
					ConnectionInfo: makeConnectionInfo(address),
				}, nil
			}
			if err.Type == k8s.ContentionError {
				tryGettingAgain = true
			}
		}
	case k8s.ContentionError:
		tryGettingAgain = true
	default:
	}
	if tryGettingAgain {
		logger.Debug("Allocate: contention detected, trying to get existing server", "matchId", request.GetMatchId())
		address, _, err = s.k8sOps.GetServerByMatchId(ctx, request.GetMatchId())
		if err == nil {
			logger.Debug("Allocate: found existing server after contention", "matchId", request.GetMatchId(), "address", address)
			return &pb.AllocateResponse{
				ConnectionInfo: makeConnectionInfo(address),
			}, nil
		}
	}
	logger.Error("Allocate: failed", "matchId", request.GetMatchId(), "errorType", err.Type, "error", err)
	return nil, convertError(err)
}

func (s *fleetManagerServiceImpl) GetServer(ctx context.Context, request *pb.GetServerRequest) (*pb.GetServerResponse, error) {
	logger := logging.GetLogger(ctx)
	if err := validateGetServerRequest(request); err != nil {
		logger.Debug("GetServer: invalid request", "error", err)
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid request: %v", err))
	}
	address, fleet, err := s.k8sOps.GetServerByMatchId(ctx, request.GetMatchId())
	if err != nil {
		if err.Type == k8s.NotFoundError {
			logger.Debug("GetServer: server not found", "matchId", request.GetMatchId())
			return &pb.GetServerResponse{}, nil
		}
		logger.Error("GetServer: failed", "matchId", request.GetMatchId(), "errorType", err.Type, "error", err)
		return nil, convertError(err)
	}
	logger.Debug("GetServer: found", "matchId", request.GetMatchId(), "address", address)
	return &pb.GetServerResponse{
		ConnectionInfo: makeConnectionInfo(address),
		Fleet:          fleet,
	}, nil
}
