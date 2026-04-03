package services

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math"
	pb "team_dynamics/api/proto/rating_service"
	pglib "team_dynamics/pg_lib/include"
	"team_dynamics/rating_service/models"
	"team_dynamics/rating_service/pg"
)

type RatingService interface {
	GetRating(ctx context.Context, request *pb.GetRatingRequest) (*pb.GetRatingResponse, error)
	UpdateRating(ctx context.Context, request *pb.UpdateRatingRequest) (*pb.UpdateRatingResponse, error)
}

type ratingServiceImpl struct {
	repo    pg.RatingRepo
	glicko2 Glicko2Service
}

func MakeRatingService(repo pg.RatingRepo, glicko2 Glicko2Service) RatingService {
	return ratingServiceImpl{repo, glicko2}
}

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

func convertMatchResult(proto pb.MatchResult) models.GameResult {
	switch proto {
	case pb.MatchResult_MATCH_RESULT_WINNER:
		return models.Winner
	case pb.MatchResult_MATCH_RESULT_LOSER:
		return models.Loser
	}
	return models.Draw
}

func makeRatingData(ratingInfo *models.RatingInfo) *pb.RatingData {
	if ratingInfo == nil {
		return nil
	}
	displayValue := int64(math.Floor(ratingInfo.Value))
	return &pb.RatingData{RatingValue: &ratingInfo.Value, DisplayValue: &displayValue}
}

func (s ratingServiceImpl) GetRating(ctx context.Context, request *pb.GetRatingRequest) (*pb.GetRatingResponse, error) {
	if request.UserId == nil {
		return nil, status.Error(codes.InvalidArgument, "no userId specified")
	}
	initialRating := s.glicko2.GetInitialRating(request.GetUserId())
	ratingInfo, err := s.repo.GetUserRating(ctx, initialRating)
	if err != nil {
		return nil, convertError(err)
	}

	return &pb.GetRatingResponse{Rating: makeRatingData(ratingInfo)}, nil
}

func (s ratingServiceImpl) UpdateRating(ctx context.Context, request *pb.UpdateRatingRequest) (*pb.UpdateRatingResponse, error) {
	if request.Player1Id == nil {
		return nil, status.Error(codes.InvalidArgument, "no player1_id specified")
	}
	if request.Player2Id == nil {
		return nil, status.Error(codes.InvalidArgument, "no player2_id specified")
	}
	if request.GetMatchResult() == pb.MatchResult_MATCH_RESULT_UNSPECIFIED {
		return nil, status.Error(codes.InvalidArgument, "no match result specified")
	}
	if request.MatchId == nil {
		return nil, status.Error(codes.InvalidArgument, "no matchId specified")
	}

	initialRating1 := s.glicko2.GetInitialRating(request.GetPlayer1Id())
	rating1Info, err := s.repo.GetUserRating(ctx, initialRating1)
	if err != nil {
		return nil, convertError(err)
	}
	initialRating2 := s.glicko2.GetInitialRating(request.GetPlayer2Id())
	rating2Info, err := s.repo.GetUserRating(ctx, initialRating2)
	if err != nil {
		return nil, convertError(err)
	}

	newRating1, newRating2 := s.glicko2.UpdateRatings(rating1Info, rating2Info, convertMatchResult(request.MatchResult))
	err = s.repo.UpdateUserRating(ctx, []*models.RatingInfo{newRating1, newRating2}, request.GetMatchId())
	if err != nil {
		return nil, convertError(err)
	}

	return &pb.UpdateRatingResponse{Player1Rating: makeRatingData(newRating1), Player2Rating: makeRatingData(newRating2)}, nil
}
