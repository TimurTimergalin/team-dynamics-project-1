package downstream

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	usPb "team_dynamics/api/proto/user_service"
)

type UserServiceClientFactory interface {
	GetUserData(ctx context.Context, req *usPb.GetUserDataRequest) (*usPb.GetUserDataResponse, error)
}

type userServiceClientFactory struct {
	address string
}

func NewUserServiceClientFactory(address string) UserServiceClientFactory {
	return &userServiceClientFactory{address}
}

func (f *userServiceClientFactory) GetUserData(ctx context.Context, req *usPb.GetUserDataRequest) (*usPb.GetUserDataResponse, error) {
	conn, err := grpc.NewClient(f.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return usPb.NewUserServiceClient(conn).GetUserData(ctx, req)
}
