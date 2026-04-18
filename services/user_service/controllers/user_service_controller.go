package controllers

import (
	"context"
	pb "team_dynamics/api/proto/user_service"
	"team_dynamics/grpc_lib"
	"team_dynamics/user_service/services"
)

type UserServiceController struct {
	pb.UnimplementedUserServiceServer
	Service services.UserService
}

func (s *UserServiceController) GetSelfData(ctx context.Context, req *pb.GetSelfDataRequest) (resp *pb.GetSelfDataResponse, err error) {
	ctx = grpc_lib.WithLogger(ctx, "GetSelfData")
	defer grpc_lib.HandlePanic(ctx, &err)
	return s.Service.GetSelfData(ctx, req)
}
func (s *UserServiceController) GetUserData(ctx context.Context, req *pb.GetUserDataRequest) (resp *pb.GetUserDataResponse, err error) {
	ctx = grpc_lib.WithLogger(ctx, "GetUserData")
	defer grpc_lib.HandlePanic(ctx, &err)
	return s.Service.GetUserData(ctx, req)
}
func (s *UserServiceController) GetFriends(ctx context.Context, req *pb.GetFriendsRequest) (resp *pb.GetFriendsResponse, err error) {
	ctx = grpc_lib.WithLogger(ctx, "GetFriends")
	defer grpc_lib.HandlePanic(ctx, &err)
	return s.Service.GetFriends(ctx, req)
}
func (s *UserServiceController) GetIncomingRequests(ctx context.Context, req *pb.GetIncomingRequestsRequest) (resp *pb.GetIncomingRequestsResponse, err error) {
	ctx = grpc_lib.WithLogger(ctx, "GetIncomingRequests")
	defer grpc_lib.HandlePanic(ctx, &err)
	return s.Service.GetIncomingRequests(ctx, req)
}
func (s *UserServiceController) GetOutgoingRequests(ctx context.Context, req *pb.GetOutgoingRequestsRequest) (resp *pb.GetOutgoingRequestsResponse, err error) {
	ctx = grpc_lib.WithLogger(ctx, "GetOutgoingRequests")
	defer grpc_lib.HandlePanic(ctx, &err)
	return s.Service.GetOutgoingRequests(ctx, req)
}

func (s *UserServiceController) AddFriend(ctx context.Context, req *pb.AddFriendRequest) (resp *pb.AddFriendResponse, err error) {
	ctx = grpc_lib.WithLogger(ctx, "GetOutgoingRequests")
	defer grpc_lib.HandlePanic(ctx, &err)
	return s.Service.AddFriend(ctx, req)
}

func (s *UserServiceController) RemoveFriend(ctx context.Context, req *pb.RemoveFriendRequest) (resp *pb.RemoveFriendResponse, err error) {
	ctx = grpc_lib.WithLogger(ctx, "GetOutgoingRequests")
	defer grpc_lib.HandlePanic(ctx, &err)
	return s.Service.RemoveFriend(ctx, req)
}
