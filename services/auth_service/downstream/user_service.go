package downstream

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	userPb "team_dynamics/api/proto/user_service"
	"team_dynamics/auth_sdk"
	"team_dynamics/logging"
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
	logger := logging.GetLogger(ctx)
	ctx, err := attachCredentials(ctx, f.authSidecar)
	if err != nil {
		logger.Error("failed to attach credentials for user service call", "error", err)
		return nil, err
	}
	conn, err := grpc.NewClient(f.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error("failed to create user service client", "address", f.address, "error", err)
		return nil, err
	}
	defer conn.Close()
	logger.Debug("calling user service GetSelfData", "address", f.address)
	resp, err := userPb.NewUserServiceClient(conn).GetSelfData(ctx, req)
	if err != nil {
		logger.Error("user service GetSelfData failed", "address", f.address, "error", err)
		return nil, err
	}
	return resp, nil
}
