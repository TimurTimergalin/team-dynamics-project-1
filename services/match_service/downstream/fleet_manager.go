package downstream

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	fmPb "team_dynamics/api/proto/fleet_manager"
)

type FleetManagerClientFactory interface {
	Allocate(ctx context.Context, req *fmPb.AllocateRequest) (*fmPb.AllocateResponse, error)
	GetServer(ctx context.Context, req *fmPb.GetServerRequest) (*fmPb.GetServerResponse, error)
}

type fleetManagerClientFactory struct {
	address string
}

func NewFleetManagerClientFactory(address string) FleetManagerClientFactory {
	return &fleetManagerClientFactory{address}
}

func (f *fleetManagerClientFactory) Allocate(ctx context.Context, req *fmPb.AllocateRequest) (*fmPb.AllocateResponse, error) {
	conn, err := grpc.NewClient(f.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return fmPb.NewFleetManagerClient(conn).Allocate(ctx, req)
}

func (f *fleetManagerClientFactory) GetServer(ctx context.Context, req *fmPb.GetServerRequest) (*fmPb.GetServerResponse, error) {
	conn, err := grpc.NewClient(f.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return fmPb.NewFleetManagerClient(conn).GetServer(ctx, req)
}
