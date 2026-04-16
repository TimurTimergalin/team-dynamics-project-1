package services

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "team_dynamics/api/proto/match_history_service"
	"team_dynamics/logging"
	"team_dynamics/match_history_service/models"
	"team_dynamics/match_history_service/pg"
	pglib "team_dynamics/pg_lib/include"
	"time"
)

type MatchHistoryService interface {
	GetMatchHistory(ctx context.Context, request *pb.GetMatchHistoryRequest) (*pb.GetMatchHistoryResponse, error)
	SaveMatch(ctx context.Context, request *pb.SaveMatchRequest) (*pb.SaveMatchResponse, error)
}

type matchHistoryServiceImpl struct {
	pageKeyService PageKeyService
	repo           pg.MatchHistoryRepo
}

func MakeMatchHistoryService(pageKeyService PageKeyService, repo pg.MatchHistoryRepo) MatchHistoryService {
	return &matchHistoryServiceImpl{pageKeyService, repo}
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

func roundFromModel(model *models.Round) *pb.Round {
	lengthMillis := model.Length.Milliseconds()
	return &pb.Round{
		IsPlayer1Killer: &model.IsPlayer1Killer,
		TimeMillis:      &lengthMillis,
	}
}

func resultFromModel(model models.MatchResult) pb.MatchResult {
	switch model {
	case models.Win:
		return pb.MatchResult_MATCH_RESULT_PLAYER1_WIN
	case models.Draw:
		return pb.MatchResult_MATCH_RESULT_DRAW
	case models.Loss:
		return pb.MatchResult_MATCH_RESULT_PLAYER2_WIN
	case models.Canceled:
		return pb.MatchResult_MATCH_RESULT_CANCELLED
	}
	return pb.MatchResult_MATCH_RESULT_UNSPECIFIED
}

func matchFromModel(model *models.AggregatedMatch) *pb.MatchData {
	rounds := make([]*pb.Round, 0, len(model.Rounds))
	for _, round := range model.Rounds {
		rounds = append(rounds, roundFromModel(round))
	}
	unixTime := model.MatchObj.End.UTC().UnixMilli()
	return &pb.MatchData{
		Player1: &pb.ParticipantData{
			Id:     &model.MatchObj.Player1Id,
			Name:   &model.MatchObj.Player1Name,
			Rating: &model.MatchObj.Player1Rating,
		},
		Player2: &pb.ParticipantData{
			Id:     &model.MatchObj.Player2Id,
			Name:   &model.MatchObj.Player2Name,
			Rating: &model.MatchObj.Player2Rating,
		},
		Rounds:       rounds,
		EndTimestamp: &unixTime,
		MatchResult:  resultFromModel(model.MatchObj.Result),
	}
}

func roundToModel(proto *pb.Round) (*models.Round, error) {
	if proto.TimeMillis == nil {
		return nil, errors.New("no time millis")
	}
	return &models.Round{
		IsPlayer1Killer: proto.GetIsPlayer1Killer(),
		Length:          time.Duration(proto.GetTimeMillis()) * time.Millisecond,
	}, nil
}

func resultToModel(proto pb.MatchResult) models.MatchResult {
	switch proto {
	case pb.MatchResult_MATCH_RESULT_PLAYER1_WIN:
		return models.Win
	case pb.MatchResult_MATCH_RESULT_DRAW:
		return models.Draw
	case pb.MatchResult_MATCH_RESULT_PLAYER2_WIN:
		return models.Loss
	}
	return models.Canceled
}

func matchToModel(proto *pb.MatchData) (*models.AggregatedMatch, error) {
	if proto.Player1 == nil {
		return nil, errors.New("no player1")
	}
	if proto.GetPlayer1().Id == nil {
		return nil, errors.New("no player1 id")
	}
	if proto.GetPlayer1().Name == nil {
		return nil, errors.New("no player1 name")
	}
	if proto.GetPlayer1().Rating == nil {
		return nil, errors.New("no player1 rating")
	}
	if proto.Player2 == nil {
		return nil, errors.New("no player1")
	}
	if proto.GetPlayer2().Id == nil {
		return nil, errors.New("no player1 id")
	}
	if proto.GetPlayer2().Name == nil {
		return nil, errors.New("no player1 name")
	}
	if proto.GetPlayer2().Rating == nil {
		return nil, errors.New("no player1 rating")
	}
	if proto.EndTimestamp == nil {
		return nil, errors.New("no end timestamp")
	}
	if proto.MatchResult == pb.MatchResult_MATCH_RESULT_UNSPECIFIED {
		return nil, errors.New("no match result")
	}
	if proto.MatchId == nil {
		return nil, errors.New("no match id")
	}

	rounds := make([]*models.Round, 0, len(proto.Rounds))
	for _, round := range proto.Rounds {
		conv, err := roundToModel(round)
		if err != nil {
			return nil, fmt.Errorf("error while parsing round: %w", err)
		}
		rounds = append(rounds, conv)
	}

	return &models.AggregatedMatch{
		MatchObj: &models.Match{
			MatchId:       proto.GetMatchId(),
			Player1Id:     proto.GetPlayer1().GetId(),
			Player1Name:   proto.GetPlayer1().GetName(),
			Player1Rating: proto.GetPlayer1().GetRating(),
			Player2Id:     proto.GetPlayer2().GetId(),
			Player2Name:   proto.GetPlayer2().GetName(),
			Player2Rating: proto.GetPlayer2().GetRating(),
			End:           time.UnixMilli(proto.GetEndTimestamp()),
			Result:        resultToModel(proto.MatchResult),
		},
		Rounds: rounds,
	}, nil
}

func (s *matchHistoryServiceImpl) GetMatchHistory(ctx context.Context, request *pb.GetMatchHistoryRequest) (*pb.GetMatchHistoryResponse, error) {
	logger := logging.GetLogger(ctx)
	if request.UserId == nil {
		logger.Debug("No user_id")
		return nil, status.Error(codes.InvalidArgument, "no user_id specified")
	}
	var matches []*models.AggregatedMatch
	var merr *pglib.PgLibError = nil
	if request.Pagekey == nil {
		logger.Debug("First page flow")
		matches, merr = s.repo.GetMatchesFirstPage(ctx, request.GetUserId())
	} else {
		logger.Debug("Second page flow")
		pageKey, err := s.pageKeyService.Deserialize(request.GetPagekey())
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid page key")
		}
		matches, merr = s.repo.GetMatchesSecondPage(ctx, request.GetUserId(), pageKey)
	}
	if merr != nil {
		return nil, convertError(merr)
	}

	matchesProto := make([]*pb.MatchData, 0, len(matches))
	for _, match := range matches {
		matchesProto = append(matchesProto, matchFromModel(match))
	}
	newPageKey := s.pageKeyService.GetNewPageKey(matches)
	return &pb.GetMatchHistoryResponse{
		Matches: matchesProto,
		Pagekey: newPageKey,
	}, nil
}

func (s *matchHistoryServiceImpl) SaveMatch(ctx context.Context, request *pb.SaveMatchRequest) (*pb.SaveMatchResponse, error) {
	if request.Match == nil {
		return nil, status.Error(codes.InvalidArgument, "no match")
	}
	aggregatedMatch, err := matchToModel(request.Match)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid match: %v", err))
	}

	pgLibError := s.repo.SaveMatch(ctx, aggregatedMatch)
	return &pb.SaveMatchResponse{}, convertError(pgLibError)
}
