package downstream

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	usPb "team_dynamics/api/proto/user_service"
	"team_dynamics/auth_sdk"
)

type UserServiceClientFactory interface {
	GetUserData(ctx context.Context, req *usPb.GetUserDataRequest) (*usPb.GetUserDataResponse, error)
}

type userServiceClientFactory struct {
	address     string
	authSidecar *auth_sdk.AuthSidecarClient
}

func NewUserServiceClientFactory(address string, authSidecar *auth_sdk.AuthSidecarClient) UserServiceClientFactory {
	return &userServiceClientFactory{address: address, authSidecar: authSidecar}
}

func (f *userServiceClientFactory) GetUserData(ctx context.Context, req *usPb.GetUserDataRequest) (*usPb.GetUserDataResponse, error) {
	ctx, err := attachCredentials(ctx, f.authSidecar)
	if err != nil {
		return nil, err
	}
	conn, err := grpc.NewClient(f.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return usPb.NewUserServiceClient(conn).GetUserData(ctx, req)
}
