package downstream

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	msPb "team_dynamics/api/proto/match_service"
)

type MatchServiceClientFactory interface {
	GetMatch(ctx context.Context, req *msPb.GetMatchRequest) (*msPb.GetMatchResponse, error)
	StartMatch(ctx context.Context, req *msPb.StartMatchRequest) (*msPb.StartMatchResponse, error)
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

func (f *matchServiceClientFactory) StartMatch(ctx context.Context, req *msPb.StartMatchRequest) (*msPb.StartMatchResponse, error) {
	conn, err := grpc.NewClient(f.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return msPb.NewMatchServiceClient(conn).StartMatch(ctx, req)
}
