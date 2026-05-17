package downstream

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	userPb "team_dynamics/api/proto/user_service"
	"team_dynamics/auth_sdk"
)

type UserServiceClientFactory interface {
	GetSelfData(ctx context.Context, req *userPb.GetSelfDataRequest) (*userPb.GetSelfDataResponse, error)
}

type userServiceClientFactory struct {
	address     string
	authSidecar *auth_sdk.AuthSidecarClient
}

func NewUserServiceClientFactory(address string, authSidecar *auth_sdk.AuthSidecarClient) UserServiceClientFactory {
	return &userServiceClientFactory{address: address, authSidecar: authSidecar}
}

func (f *userServiceClientFactory) GetSelfData(ctx context.Context, req *userPb.GetSelfDataRequest) (*userPb.GetSelfDataResponse, error) {
	ctx, err := attachCredentials(ctx, f.authSidecar)
	if err != nil {
		return nil, err
	}
	conn, err := grpc.NewClient(f.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return userPb.NewUserServiceClient(conn).GetSelfData(ctx, req)
}
