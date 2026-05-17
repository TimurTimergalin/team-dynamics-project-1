package downstream

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	rsPb "team_dynamics/api/proto/rating_service"
	"team_dynamics/auth_sdk"
)

type RatingServiceClientFactory interface {
	GetRating(ctx context.Context, req *rsPb.GetRatingRequest) (*rsPb.GetRatingResponse, error)
}

type ratingServiceClientFactory struct {
	address     string
	authSidecar *auth_sdk.AuthSidecarClient
}

func NewRatingServiceClientFactory(address string, authSidecar *auth_sdk.AuthSidecarClient) RatingServiceClientFactory {
	return &ratingServiceClientFactory{address: address, authSidecar: authSidecar}
}

func (f *ratingServiceClientFactory) GetRating(ctx context.Context, req *rsPb.GetRatingRequest) (*rsPb.GetRatingResponse, error) {
	ctx, err := attachCredentials(ctx, f.authSidecar)
	if err != nil {
		return nil, err
	}
	conn, err := grpc.NewClient(f.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return rsPb.NewRatingServiceClient(conn).GetRating(ctx, req)
}
