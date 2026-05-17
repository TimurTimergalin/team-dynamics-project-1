package downstream

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	authPb "team_dynamics/api/proto/auth_service"
	"team_dynamics/auth_sdk"
)

type AuthServiceClient interface {
	GetPublicKey(ctx context.Context, req *authPb.GetPublicKeyRequest) (*authPb.GetPublicKeyResponse, error)
}

type authServiceClient struct {
	address     string
	authSidecar *auth_sdk.AuthSidecarClient
}

func NewAuthServiceClient(address string, authSidecar *auth_sdk.AuthSidecarClient) AuthServiceClient {
	return &authServiceClient{address: address, authSidecar: authSidecar}
}

func (c *authServiceClient) GetPublicKey(ctx context.Context, req *authPb.GetPublicKeyRequest) (*authPb.GetPublicKeyResponse, error) {
	ctx, err := attachCredentials(ctx, c.authSidecar)
	if err != nil {
		return nil, err
	}
	conn, err := grpc.NewClient(c.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return authPb.NewAuthServiceClient(conn).GetPublicKey(ctx, req)
}
