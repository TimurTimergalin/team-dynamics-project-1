package downstream

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	mhsPb "team_dynamics/api/proto/match_history_service"
)

type MatchHistoryServiceClientFactory interface {
	SaveMatch(ctx context.Context, req *mhsPb.SaveMatchRequest) (*mhsPb.SaveMatchResponse, error)
}

type matchHistoryServiceClientFactory struct {
	address string
}

func NewMatchHistoryServiceClientFactory(address string) MatchHistoryServiceClientFactory {
	return &matchHistoryServiceClientFactory{address}
}

func (f *matchHistoryServiceClientFactory) SaveMatch(ctx context.Context, req *mhsPb.SaveMatchRequest) (*mhsPb.SaveMatchResponse, error) {
	conn, err := grpc.NewClient(f.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return mhsPb.NewMatchHistoryServiceClient(conn).SaveMatch(ctx, req)
}
