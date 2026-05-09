package services

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pbCommon "team_dynamics/api/proto/user_common"
	pb "team_dynamics/api/proto/user_service"
	"team_dynamics/logging"
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
	repo           pg.UserStorageRepo
	pageKeyService PageKeyService
	steamService   SteamService
}

func MakeUserService(repo pg.UserStorageRepo, pageKeyService PageKeyService, steamService SteamService) UserService {
	return &userServiceImpl{repo, pageKeyService, steamService}
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

func convertFriend(f *models.Friend) *pb.Friend {
	if f == nil {
		return nil
	}
	source := fmt.Sprintf("%d", f.Source)
	return &pb.Friend{
		User:   convertUserData(f.Data),
		Source: &source,
	}
}

func convertFriends(friends []*models.Friend) []*pb.Friend {
	result := make([]*pb.Friend, 0, len(friends))
	for _, f := range friends {
		result = append(result, convertFriend(f))
	}
	return result
}

func validateGetFriendsRequest(req *pb.GetFriendsRequest) error {
	if req == nil || req.UserId == nil {
		return errors.New("user_id is not set")
	}
	return nil
}

func validateGetIncomingRequestsRequest(req *pb.GetIncomingRequestsRequest) error {
	if req == nil || req.UserId == nil {
		return errors.New("user_id is not set")
	}
	return nil
}

func validateGetOutgoingRequestsRequest(req *pb.GetOutgoingRequestsRequest) error {
	if req == nil || req.UserId == nil {
		return errors.New("user_id is not set")
	}
	return nil
}

func (s *userServiceImpl) parsePageKey(raw *string) (*models.PageKey, error) {
	if raw == nil {
		return &models.PageKey{LastUserId: 0}, nil
	}
	return s.pageKeyService.Deserialize(*raw)
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
	logger := logging.GetLogger(ctx)
	if err := validateGetSelfDataRequest(req); err != nil {
		logger.Debug("invalid request", "error", err)
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid request: %v", err))
	}
	steamId := req.Key.Key.(*pbCommon.ExternalKey_SteamId).SteamId
	steamData, err := s.steamService.GetUserSummary(ctx, fmt.Sprintf("%d", steamId))
	if err != nil {
		logger.Error("failed to fetch Steam data", "steam_id", steamId, "error", err)
		return nil, status.Errorf(codes.Unavailable, "failed to fetch Steam data: %v", err)
	}
	user, pgErr := s.repo.UpsertSelfData(ctx, steamId, steamData.Name)
	if pgErr != nil {
		logger.Error("failed to upsert user", "steam_id", steamId, "error", pgErr)
		return nil, convertError(pgErr)
	}
	logger.Debug("GetSelfData: ok", "user_id", user.Id)
	return &pb.GetSelfDataResponse{
		UserData: convertUserData(user),
	}, nil
}

func (s *userServiceImpl) GetUserData(ctx context.Context, req *pb.GetUserDataRequest) (*pb.GetUserDataResponse, error) {
	logger := logging.GetLogger(ctx)
	if err := validateGetUserDataRequest(req); err != nil {
		logger.Debug("invalid request", "error", err)
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid request: %v", err))
	}
	userId := *req.Id
	user, err := s.repo.GetUserData(ctx, userId)
	if err != nil {
		logger.Error("failed to get user data", "user_id", userId, "error", err)
		return nil, convertError(err)
	}
	logger.Debug("GetUserData: ok", "user_id", userId)
	return &pb.GetUserDataResponse{
		UserData: convertUserData(user),
	}, nil
}

func (s *userServiceImpl) GetFriends(ctx context.Context, req *pb.GetFriendsRequest) (*pb.GetFriendsResponse, error) {
	logger := logging.GetLogger(ctx)
	if err := validateGetFriendsRequest(req); err != nil {
		logger.Debug("invalid request", "error", err)
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid request: %v", err))
	}
	pageKey, err := s.parsePageKey(req.Pagekey)
	if err != nil {
		logger.Debug("invalid pagekey", "error", err)
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid pagekey: %v", err))
	}
	friends, pgErr := s.repo.GetFriends(ctx, *req.UserId, pageKey)
	if pgErr != nil {
		logger.Error("failed to get friends", "user_id", *req.UserId, "error", pgErr)
		return nil, convertError(pgErr)
	}
	logger.Debug("GetFriends: ok", "user_id", *req.UserId, "count", len(friends))
	return &pb.GetFriendsResponse{
		Friends: convertFriends(friends),
		Pagekey: s.pageKeyService.GetNewPageKey(friends),
	}, nil
}

func (s *userServiceImpl) GetIncomingRequests(ctx context.Context, req *pb.GetIncomingRequestsRequest) (*pb.GetIncomingRequestsResponse, error) {
	logger := logging.GetLogger(ctx)
	if err := validateGetIncomingRequestsRequest(req); err != nil {
		logger.Debug("invalid request", "error", err)
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid request: %v", err))
	}
	pageKey, err := s.parsePageKey(req.Pagekey)
	if err != nil {
		logger.Debug("invalid pagekey", "error", err)
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid pagekey: %v", err))
	}
	friends, pgErr := s.repo.GetIncomingRequests(ctx, *req.UserId, pageKey)
	if pgErr != nil {
		logger.Error("failed to get incoming requests", "user_id", *req.UserId, "error", pgErr)
		return nil, convertError(pgErr)
	}
	logger.Debug("GetIncomingRequests: ok", "user_id", *req.UserId, "count", len(friends))
	return &pb.GetIncomingRequestsResponse{
		Friends: convertFriends(friends),
		Pagekey: s.pageKeyService.GetNewPageKey(friends),
	}, nil
}

func (s *userServiceImpl) GetOutgoingRequests(ctx context.Context, req *pb.GetOutgoingRequestsRequest) (*pb.GetOutgoingRequestsResponse, error) {
	logger := logging.GetLogger(ctx)
	if err := validateGetOutgoingRequestsRequest(req); err != nil {
		logger.Debug("invalid request", "error", err)
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid request: %v", err))
	}
	pageKey, err := s.parsePageKey(req.Pagekey)
	if err != nil {
		logger.Debug("invalid pagekey", "error", err)
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid pagekey: %v", err))
	}
	friends, pgErr := s.repo.GetOutgoingRequests(ctx, *req.UserId, pageKey)
	if pgErr != nil {
		logger.Error("failed to get outgoing requests", "user_id", *req.UserId, "error", pgErr)
		return nil, convertError(pgErr)
	}
	logger.Debug("GetOutgoingRequests: ok", "user_id", *req.UserId, "count", len(friends))
	return &pb.GetOutgoingRequestsResponse{
		Friends: convertFriends(friends),
		Pagekey: s.pageKeyService.GetNewPageKey(friends),
	}, nil
}

func validateAddFriendRequest(req *pb.AddFriendRequest) error {
	if req == nil || req.UserId == nil {
		return errors.New("user_id is not set")
	}
	if req.OtherUserId == nil {
		return errors.New("other_user_id is not set")
	}
	if *req.UserId == *req.OtherUserId {
		return errors.New("user_id and other_user_id must be different")
	}
	return nil
}

func validateRemoveFriendRequest(req *pb.RemoveFriendRequest) error {
	if req == nil || req.UserId == nil {
		return errors.New("user_id is not set")
	}
	if req.OtherUserId == nil {
		return errors.New("other_user_id is not set")
	}
	if *req.UserId == *req.OtherUserId {
		return errors.New("user_id and other_user_id must be different")
	}
	return nil
}

func convertAddFriendResult(r models.AddFriendResult) pb.AddFriendResult {
	switch r {
	case models.AddFriendRequestSent:
		return pb.AddFriendResult_ADD_FRIEND_RESULT_REQUEST_SENT
	case models.AddFriendAccepted:
		return pb.AddFriendResult_ADD_FRIEND_RESULT_ACCEPTED
	default:
		return pb.AddFriendResult_ADD_FRIEND_RESULT_NOOP
	}
}

func convertRemoveFriendResult(r models.RemoveFriendResult) pb.RemoveFriendResult {
	switch r {
	case models.RemoveFriendRequestCancelled:
		return pb.RemoveFriendResult_REMOVE_FRIEND_RESULT_REQUEST_CANCELLED
	case models.RemoveFriendRequestDeclined:
		return pb.RemoveFriendResult_REMOVE_FRIEND_RESULT_REQUEST_DECLINED
	case models.RemoveFriendFriendRemoved:
		return pb.RemoveFriendResult_REMOVE_FRIEND_RESULT_FRIEND_REMOVED
	default:
		return pb.RemoveFriendResult_REMOVE_FRIEND_RESULT_NOOP
	}
}

func (s *userServiceImpl) AddFriend(ctx context.Context, req *pb.AddFriendRequest) (*pb.AddFriendResponse, error) {
	logger := logging.GetLogger(ctx)
	if err := validateAddFriendRequest(req); err != nil {
		logger.Debug("invalid request", "error", err)
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid request: %v", err))
	}
	friend, result, pgErr := s.repo.AddFriend(ctx, *req.UserId, *req.OtherUserId)
	if pgErr != nil {
		logger.Error("failed to add friend", "user_id", *req.UserId, "other_user_id", *req.OtherUserId, "error", pgErr)
		return nil, convertError(pgErr)
	}
	logger.Debug("AddFriend: ok", "user_id", *req.UserId, "other_user_id", *req.OtherUserId, "result", result)
	return &pb.AddFriendResponse{
		Friend: convertFriend(friend),
		Result: convertAddFriendResult(result),
	}, nil
}

func (s *userServiceImpl) RemoveFriend(ctx context.Context, req *pb.RemoveFriendRequest) (*pb.RemoveFriendResponse, error) {
	logger := logging.GetLogger(ctx)
	if err := validateRemoveFriendRequest(req); err != nil {
		logger.Debug("invalid request", "error", err)
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid request: %v", err))
	}
	result, pgErr := s.repo.RemoveFriend(ctx, *req.UserId, *req.OtherUserId)
	if pgErr != nil {
		logger.Error("failed to remove friend", "user_id", *req.UserId, "other_user_id", *req.OtherUserId, "error", pgErr)
		return nil, convertError(pgErr)
	}
	logger.Debug("RemoveFriend: ok", "user_id", *req.UserId, "other_user_id", *req.OtherUserId, "result", result)
	return &pb.RemoveFriendResponse{
		Result: convertRemoveFriendResult(result),
	}, nil
}
