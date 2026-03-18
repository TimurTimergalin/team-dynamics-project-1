package controllers

import (
	"context"
	pb "fleet_manager/fleet_manager/proto"
	"fleet_manager/services"
)

type FleetManagerController struct {
	pb.UnimplementedFleetManagerServer
	Service *services.FleetManagerService
}

func (s *FleetManagerController) Allocate(ctx context.Context, req *pb.AllocateRequest) (*pb.AllocateResponse, error) {
	return s.Service.Allocate(ctx, req)
}

func (s *FleetManagerController) GetServer(ctx context.Context, req *pb.GetServerRequest) (*pb.GetServerResponse, error) {
	return s.Service.GetServer(ctx, req)
}
