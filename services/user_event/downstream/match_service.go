package downstream

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	msPb "team_dynamics/api/proto/match_service"
	"team_dynamics/auth_sdk"
)

type MatchServiceClientFactory interface {
	GetMatch(ctx context.Context, req *msPb.GetMatchRequest) (*msPb.GetMatchResponse, error)
	StartMatch(ctx context.Context, req *msPb.StartMatchRequest) (*msPb.StartMatchResponse, error)
}

type matchServiceClientFactory struct {
	address     string
	authSidecar *auth_sdk.AuthSidecarClient
}

func NewMatchServiceClientFactory(address string, authSidecar *auth_sdk.AuthSidecarClient) MatchServiceClientFactory {
	return &matchServiceClientFactory{address: address, authSidecar: authSidecar}
}

func (f *matchServiceClientFactory) GetMatch(ctx context.Context, req *msPb.GetMatchRequest) (*msPb.GetMatchResponse, error) {
	ctx, err := attachCredentials(ctx, f.authSidecar)
	if err != nil {
		return nil, err
	}
	conn, err := grpc.NewClient(f.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return msPb.NewMatchServiceClient(conn).GetMatch(ctx, req)
}

func (f *matchServiceClientFactory) StartMatch(ctx context.Context, req *msPb.StartMatchRequest) (*msPb.StartMatchResponse, error) {
	ctx, err := attachCredentials(ctx, f.authSidecar)
	if err != nil {
		return nil, err
	}
	conn, err := grpc.NewClient(f.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return msPb.NewMatchServiceClient(conn).StartMatch(ctx, req)
}
