package downstream

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	msPb "team_dynamics/api/proto/match_service"
)

type MatchServiceClientFactory interface {
	GetMatch(ctx context.Context, req *msPb.GetMatchRequest) (*msPb.GetMatchResponse, error)
}

type matchServiceClientFactory struct {
	address string
}

func NewMatchServiceClientFactory(address string) MatchServiceClientFactory {
	return &matchServiceClientFactory{address}
}

func (f *matchServiceClientFactory) GetMatch(ctx context.Context, req *msPb.GetMatchRequest) (*msPb.GetMatchResponse, error) {
	conn, err := grpc.NewClient(f.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return msPb.NewMatchServiceClient(conn).GetMatch(ctx, req)
}
