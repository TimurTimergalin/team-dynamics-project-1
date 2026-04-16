package controllers

import (
	"context"
	pb "team_dynamics/api/proto/match_service"
	"team_dynamics/grpc_lib"
	"team_dynamics/match_service/services"
)

type MatchServiceController struct {
	pb.UnimplementedMatchServiceServer
	Service services.MatchService
}

func (c *MatchServiceController) StartMatch(ctx context.Context, request *pb.StartMatchRequest) (response *pb.StartMatchResponse, err error) {
	ctx = grpc_lib.WithLogger(ctx, "StartMatch")
	defer grpc_lib.HandlePanic(ctx, &err)
	return c.Service.StartMatch(ctx, request)
}
func (c *MatchServiceController) GetMatch(ctx context.Context, request *pb.GetMatchRequest) (response *pb.GetMatchResponse, err error) {
	ctx = grpc_lib.WithLogger(ctx, "GetMatch")
	defer grpc_lib.HandlePanic(ctx, &err)
	return c.Service.GetMatch(ctx, request)
}
func (c *MatchServiceController) EndMatch(ctx context.Context, request *pb.EndMatchRequest) (response *pb.EndMatchResponse, err error) {
	ctx = grpc_lib.WithLogger(ctx, "EndMatch")
	defer grpc_lib.HandlePanic(ctx, &err)
	return c.Service.EndMatch(ctx, request)
}
func (c *MatchServiceController) CancelMatch(ctx context.Context, request *pb.CancelMatchRequest) (response *pb.CancelMatchResponse, err error) {
	ctx = grpc_lib.WithLogger(ctx, "CancelMatch")
	defer grpc_lib.HandlePanic(ctx, &err)
	return c.Service.CancelMatch(ctx, request)
}
func (c *MatchServiceController) RenewMatch(ctx context.Context, request *pb.RenewMatchRequest) (response *pb.RenewMatchResponse, err error) {
	ctx = grpc_lib.WithLogger(ctx, "RenewMatch")
	defer grpc_lib.HandlePanic(ctx, &err)
	return c.Service.RenewMatch(ctx, request)
}
