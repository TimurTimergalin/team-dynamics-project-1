package controllers

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (c *RatingServiceController) GetRating(ctx context.Context, request *pb.GetRatingRequest) (response *pb.GetRatingResponse, err error) {
	defer handlePanic(&err)
	response, err = c.Service.GetRating(ctx, request)
	return
}

func (c *RatingServiceController) UpdateRating(ctx context.Context, request *pb.UpdateRatingRequest) (response *pb.UpdateRatingResponse, err error) {
	defer handlePanic(&err)
	response, err = c.Service.UpdateRating(ctx, request)
	return
}
