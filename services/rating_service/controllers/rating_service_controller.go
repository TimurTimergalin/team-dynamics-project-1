package controllers

import (
	"context"
	pb "team_dynamics/api/proto/rating_service"
	"team_dynamics/grpc_lib"
	"team_dynamics/rating_service/services"
)

type RatingServiceController struct {
	pb.UnimplementedRatingServiceServer
	Service services.RatingService
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
