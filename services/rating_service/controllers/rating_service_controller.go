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
	defer handlePanic(&err)
	response, err = c.Service.GetRating(addLogger(ctx, "GetRating"), request)
	return
}

func (c *RatingServiceController) UpdateRating(ctx context.Context, request *pb.UpdateRatingRequest) (response *pb.UpdateRatingResponse, err error) {
	defer handlePanic(&err)
	response, err = c.Service.UpdateRating(addLogger(ctx, "UpdateRating"), request)
	return
}
