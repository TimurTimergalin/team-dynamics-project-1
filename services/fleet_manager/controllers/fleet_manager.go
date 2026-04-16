package controllers

import (
	"context"
	pb "team_dynamics/api/proto/fleet_manager"
	"team_dynamics/fleet_manager/services"
	"team_dynamics/grpc_lib"
)

type FleetManagerController struct {
	pb.UnimplementedFleetManagerServer
	Service services.FleetManagerService
}

func (s *FleetManagerController) Allocate(ctx context.Context, req *pb.AllocateRequest) (resp *pb.AllocateResponse, err error) {
	ctx = grpc_lib.WithLogger(ctx, "Allocate")
	defer grpc_lib.HandlePanic(ctx, &err)
	return s.Service.Allocate(ctx, req)
}

func (s *FleetManagerController) GetServer(ctx context.Context, req *pb.GetServerRequest) (resp *pb.GetServerResponse, err error) {
	ctx = grpc_lib.WithLogger(ctx, "GetServer")
	defer grpc_lib.HandlePanic(ctx, &err)
	return s.Service.GetServer(ctx, req)
}
