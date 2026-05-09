package downstream

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	authPb "team_dynamics/api/proto/auth_service"
)

type AuthServiceClient interface {
	GetPublicKey(ctx context.Context, req *authPb.GetPublicKeyRequest) (*authPb.GetPublicKeyResponse, error)
}

type authServiceClient struct {
	address string
}

func NewAuthServiceClient(address string) AuthServiceClient {
	return &authServiceClient{address}
}

func (c *authServiceClient) GetPublicKey(ctx context.Context, req *authPb.GetPublicKeyRequest) (*authPb.GetPublicKeyResponse, error) {
	conn, err := grpc.NewClient(c.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return authPb.NewAuthServiceClient(conn).GetPublicKey(ctx, req)
}
