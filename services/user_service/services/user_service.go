package services

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pbCommon "team_dynamics/api/proto/user_common"
	pb "team_dynamics/api/proto/user_service"
	pglib "team_dynamics/pg_lib/include"
	"team_dynamics/user_service/models"
	"team_dynamics/user_service/pg"
)

func convertError(pgLibErr *pglib.PgLibError) error {
	if pgLibErr == nil {
		return nil
	}
	switch pgLibErr.Type {
	case pglib.LogicError:
		return status.Errorf(codes.FailedPrecondition, "logic error encountered: %v", pgLibErr)
	case pglib.ServerError:
		return status.Errorf(codes.Internal, "server error encountered: %v", pgLibErr)
	case pglib.ConnectionError:
		return status.Errorf(codes.Canceled, "connection error encountered: %v", pgLibErr)
	}
	return pgLibErr
}

type UserService interface {
	GetSelfData(ctx context.Context, req *pb.GetSelfDataRequest) (*pb.GetSelfDataResponse, error)
	GetUserData(ctx context.Context, req *pb.GetUserDataRequest) (*pb.GetUserDataResponse, error)
	GetFriends(ctx context.Context, req *pb.GetFriendsRequest) (*pb.GetFriendsResponse, error)
	GetIncomingRequests(ctx context.Context, req *pb.GetIncomingRequestsRequest) (*pb.GetIncomingRequestsResponse, error)
	GetOutgoingRequests(ctx context.Context, req *pb.GetOutgoingRequestsRequest) (*pb.GetOutgoingRequestsResponse, error)
	AddFriend(ctx context.Context, req *pb.AddFriendRequest) (*pb.AddFriendResponse, error)
	RemoveFriend(ctx context.Context, req *pb.RemoveFriendRequest) (*pb.RemoveFriendResponse, error)
}

type userServiceImpl struct {
	repo pg.UserStorageRepo
}

func MakeUserService(repo pg.UserStorageRepo) UserService {
	return &userServiceImpl{repo}
}

func convertUserData(model *models.UserData) *pb.UserData {
	if model == nil {
		return nil
	}
	return &pb.UserData{
		Id:   &model.Id,
		Name: &model.Name,
	}
}

func validateGetSelfDataRequest(req *pb.GetSelfDataRequest) error {
	if req == nil || req.Key == nil || req.Key.Key == nil {
		return errors.New("steam id is not set")
	}
	return nil
}

func validateGetUserDataRequest(req *pb.GetUserDataRequest) error {
	if req == nil || req.Id == nil {
		return errors.New("id is not set")
	}
	return nil
}

func (s *userServiceImpl) GetSelfData(ctx context.Context, req *pb.GetSelfDataRequest) (*pb.GetSelfDataResponse, error) {
	if err := validateGetSelfDataRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid request: %v", err))
	}
	steamId := req.Key.Key.(*pbCommon.ExternalKey_SteamId).SteamId
	user, err := s.repo.GetSelfData(ctx, steamId)
	if err != nil {
		return nil, convertError(err)
	}
	return &pb.GetSelfDataResponse{
		UserData: convertUserData(user),
	}, nil
}

func (s *userServiceImpl) GetUserData(ctx context.Context, req *pb.GetUserDataRequest) (*pb.GetUserDataResponse, error) {
	if err := validateGetUserDataRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid request: %v", err))
	}
	userId := *req.Id
	user, err := s.repo.GetUserData(ctx, userId)
	if err != nil {
		return nil, convertError(err)
	}
	return &pb.GetUserDataResponse{
		UserData: convertUserData(user),
	}, nil
}

func (s *userServiceImpl) GetFriends(ctx context.Context, req *pb.GetFriendsRequest) (*pb.GetFriendsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *userServiceImpl) GetIncomingRequests(ctx context.Context, req *pb.GetIncomingRequestsRequest) (*pb.GetIncomingRequestsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *userServiceImpl) GetOutgoingRequests(ctx context.Context, req *pb.GetOutgoingRequestsRequest) (*pb.GetOutgoingRequestsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *userServiceImpl) AddFriend(ctx context.Context, req *pb.AddFriendRequest) (*pb.AddFriendResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *userServiceImpl) RemoveFriend(ctx context.Context, req *pb.RemoveFriendRequest) (*pb.RemoveFriendResponse, error) {
	//TODO implement me
	panic("implement me")
}
