package downstream

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "team_dynamics/api/proto/match_history_service_v2"
)

type MatchHistoryServiceV2ClientFactory interface {
	GetRating(ctx context.Context, req *pb.GetRatingRequest) (*pb.GetRatingResponse, error)
}

type matchHistoryServiceV2ClientFactory struct {
	address string
}

func NewMatchHistoryServiceV2ClientFactory(address string) MatchHistoryServiceV2ClientFactory {
	return &matchHistoryServiceV2ClientFactory{address}
}

func (f *matchHistoryServiceV2ClientFactory) GetRating(ctx context.Context, req *pb.GetRatingRequest) (*pb.GetRatingResponse, error) {
	conn, err := grpc.NewClient(f.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return pb.NewMatchHistoryServiceV2Client(conn).GetRating(ctx, req)
}
