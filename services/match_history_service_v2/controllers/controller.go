package controllers

import (
	"context"
	pb "team_dynamics/api/proto/match_history_service_v2"
	"team_dynamics/grpc_lib"
	"team_dynamics/match_history_service_v2/services"
)

type MatchHistoryServiceV2Controller struct {
	pb.UnimplementedMatchHistoryServiceV2Server
	Service services.MatchHistoryServiceV2
}

func (c *MatchHistoryServiceV2Controller) GetMatchHistory(ctx context.Context, req *pb.GetMatchHistoryRequest) (resp *pb.GetMatchHistoryResponse, err error) {
	ctx = grpc_lib.WithLogger(ctx, "GetMatchHistory")
	defer grpc_lib.HandlePanic(ctx, &err)
	return c.Service.GetMatchHistory(ctx, req)
}

func (c *MatchHistoryServiceV2Controller) GetRating(ctx context.Context, req *pb.GetRatingRequest) (resp *pb.GetRatingResponse, err error) {
	ctx = grpc_lib.WithLogger(ctx, "GetRating")
	defer grpc_lib.HandlePanic(ctx, &err)
	return c.Service.GetRating(ctx, req)
}

func (c *MatchHistoryServiceV2Controller) SaveMatch(ctx context.Context, req *pb.SaveMatchRequest) (resp *pb.SaveMatchResponse, err error) {
	ctx = grpc_lib.WithLogger(ctx, "SaveMatch")
	defer grpc_lib.HandlePanic(ctx, &err)
	return c.Service.SaveMatch(ctx, req)
}
