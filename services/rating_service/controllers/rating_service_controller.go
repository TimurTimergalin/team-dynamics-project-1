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
	pb "team_dynamics/api/proto/rating_service"
	"team_dynamics/grpc_lib"
	"team_dynamics/rating_service/services"
)

type RatingServiceController struct {
	pb.UnimplementedRatingServiceServer
	Service services.RatingService
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

func (c *RatingServiceController) GetRating(ctx context.Context, request *pb.GetRatingRequest) (response *pb.GetRatingResponse, err error) {
	ctx = grpc_lib.WithLogger(ctx, "GetRating")
	defer grpc_lib.HandlePanic(ctx, &err)
	return c.Service.GetRating(ctx, request)
}

func (c *RatingServiceController) UpdateRating(ctx context.Context, request *pb.UpdateRatingRequest) (response *pb.UpdateRatingResponse, err error) {
	ctx = grpc_lib.WithLogger(ctx, "UpdateRating")
	defer grpc_lib.HandlePanic(ctx, &err)
	return c.Service.UpdateRating(ctx, request)
}
