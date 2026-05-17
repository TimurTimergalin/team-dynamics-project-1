package downstream

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	fmPb "team_dynamics/api/proto/fleet_manager"
	"team_dynamics/auth_sdk"
)

type FleetManagerClientFactory interface {
	Allocate(ctx context.Context, req *fmPb.AllocateRequest) (*fmPb.AllocateResponse, error)
	GetServer(ctx context.Context, req *fmPb.GetServerRequest) (*fmPb.GetServerResponse, error)
}

type fleetManagerClientFactory struct {
	address     string
	authSidecar *auth_sdk.AuthSidecarClient
}

func NewFleetManagerClientFactory(address string, authSidecar *auth_sdk.AuthSidecarClient) FleetManagerClientFactory {
	return &fleetManagerClientFactory{address: address, authSidecar: authSidecar}
}

func (f *fleetManagerClientFactory) Allocate(ctx context.Context, req *fmPb.AllocateRequest) (*fmPb.AllocateResponse, error) {
	ctx, err := attachCredentials(ctx, f.authSidecar)
	if err != nil {
		return nil, err
	}
	conn, err := grpc.NewClient(f.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return fmPb.NewFleetManagerClient(conn).Allocate(ctx, req)
}

func (f *fleetManagerClientFactory) GetServer(ctx context.Context, req *fmPb.GetServerRequest) (*fmPb.GetServerResponse, error) {
	ctx, err := attachCredentials(ctx, f.authSidecar)
	if err != nil {
		return nil, err
	}
	conn, err := grpc.NewClient(f.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return fmPb.NewFleetManagerClient(conn).GetServer(ctx, req)
}
