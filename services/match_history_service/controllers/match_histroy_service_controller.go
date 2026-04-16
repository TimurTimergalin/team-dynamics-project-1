package controllers

import (
	"context"
	pb "team_dynamics/api/proto/match_history_service"
	"team_dynamics/grpc_lib"
	"team_dynamics/match_history_service/services"
)

type MatchHistoryServiceController struct {
	pb.UnimplementedMatchHistoryServiceServer
	Service services.MatchHistoryService
}

func (c *MatchHistoryServiceController) GetMatchHistory(ctx context.Context, request *pb.GetMatchHistoryRequest) (response *pb.GetMatchHistoryResponse, err error) {
	ctx = grpc_lib.WithLogger(ctx, "GetMatchHistory")
	defer grpc_lib.HandlePanic(ctx, &err)
	return c.Service.GetMatchHistory(ctx, request)
}

func (c *MatchHistoryServiceController) SaveMatch(ctx context.Context, request *pb.SaveMatchRequest) (response *pb.SaveMatchResponse, err error) {
	ctx = grpc_lib.WithLogger(ctx, "SaveMatch")
	defer grpc_lib.HandlePanic(ctx, &err)
	return c.Service.SaveMatch(ctx, request)
}
