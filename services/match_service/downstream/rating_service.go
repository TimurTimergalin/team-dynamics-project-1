package downstream

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	rsPb "team_dynamics/api/proto/rating_service"
)

type RatingServiceClientFactory interface {
	UpdateRating(ctx context.Context, req *rsPb.UpdateRatingRequest) (*rsPb.UpdateRatingResponse, error)
}

type ratingServiceClientFactory struct {
	address string
}

func NewRatingServiceClientFactory(address string) RatingServiceClientFactory {
	return &ratingServiceClientFactory{address}
}

func (f *ratingServiceClientFactory) UpdateRating(ctx context.Context, req *rsPb.UpdateRatingRequest) (*rsPb.UpdateRatingResponse, error) {
	conn, err := grpc.NewClient(f.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return rsPb.NewRatingServiceClient(conn).UpdateRating(ctx, req)
}
