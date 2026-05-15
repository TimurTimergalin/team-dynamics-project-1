package downstream

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	userPb "team_dynamics/api/proto/user_service"
)

type UserServiceClientFactory interface {
	GetSelfData(ctx context.Context, req *userPb.GetSelfDataRequest) (*userPb.GetSelfDataResponse, error)
}

type userServiceClientFactory struct {
	address string
}

func NewUserServiceClientFactory(address string) UserServiceClientFactory {
	return &userServiceClientFactory{address}
}

func (f *userServiceClientFactory) GetSelfData(ctx context.Context, req *userPb.GetSelfDataRequest) (*userPb.GetSelfDataResponse, error) {
	conn, err := grpc.NewClient(f.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return userPb.NewUserServiceClient(conn).GetSelfData(ctx, req)
}