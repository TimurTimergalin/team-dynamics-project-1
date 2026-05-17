package downstream

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	mhsPb "team_dynamics/api/proto/match_history_service"
	"team_dynamics/auth_sdk"
)

type MatchHistoryServiceClientFactory interface {
	SaveMatch(ctx context.Context, req *mhsPb.SaveMatchRequest) (*mhsPb.SaveMatchResponse, error)
}

type matchHistoryServiceClientFactory struct {
	address     string
	authSidecar *auth_sdk.AuthSidecarClient
}

func NewMatchHistoryServiceClientFactory(address string, authSidecar *auth_sdk.AuthSidecarClient) MatchHistoryServiceClientFactory {
	return &matchHistoryServiceClientFactory{address: address, authSidecar: authSidecar}
}

func (f *matchHistoryServiceClientFactory) SaveMatch(ctx context.Context, req *mhsPb.SaveMatchRequest) (*mhsPb.SaveMatchResponse, error) {
	ctx, err := attachCredentials(ctx, f.authSidecar)
	if err != nil {
		return nil, err
	}
	conn, err := grpc.NewClient(f.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return mhsPb.NewMatchHistoryServiceClient(conn).SaveMatch(ctx, req)
}
