package controllers

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"logging"
	"os"
	pb "team_dynamics/api/proto/match_history_service"
	"team_dynamics/match_history_service/services"
)

type MatchHistoryServiceController struct {
	pb.UnimplementedMatchHistoryServiceServer
	Service services.MatchHistoryService
}

func handlePanic(errOut *error) {
	if r := recover(); r != nil {
		*errOut = status.Error(codes.Internal, fmt.Sprintf("panic occured: %v", r))
	}
}

func addLogger(ctx context.Context, rpc string) context.Context {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug})
	logger := slog.New(handler).With("requestId", uuid.New().String(), "rpc", rpc)
	return logging.WithLogger(ctx, logger)
}

func (c *MatchHistoryServiceController) GetMatchHistory(ctx context.Context, request *pb.GetMatchHistoryRequest) (response *pb.GetMatchHistoryResponse, err error) {
	defer handlePanic(&err)
	return c.Service.GetMatchHistory(addLogger(ctx, "GetMatchHistory"), request)
}

func (c *MatchHistoryServiceController) SaveMatch(ctx context.Context, request *pb.SaveMatchRequest) (response *pb.SaveMatchResponse, err error) {
	defer handlePanic(&err)
	return c.Service.SaveMatch(addLogger(ctx, "SaveMatch"), request)
}
