package services

import (
	"context"
	"errors"
	"fmt"
	"math"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "team_dynamics/api/proto/match_history_service_v2"
	"team_dynamics/logging"
	"team_dynamics/match_history_service_v2/models"
	"team_dynamics/match_history_service_v2/pg"
	pglib "team_dynamics/pg_lib/include"
	"time"
)

type MatchHistoryServiceV2 interface {
	GetMatchHistory(ctx context.Context, req *pb.GetMatchHistoryRequest) (*pb.GetMatchHistoryResponse, error)
	GetRating(ctx context.Context, req *pb.GetRatingRequest) (*pb.GetRatingResponse, error)
	SaveMatch(ctx context.Context, req *pb.SaveMatchRequest) (*pb.SaveMatchResponse, error)
}

type matchHistoryServiceV2Impl struct {
	repo           pg.MatchHistoryV2Repo
	pageKeyService PageKeyService
	glicko2        Glicko2Service
}

func MakeMatchHistoryServiceV2(repo pg.MatchHistoryV2Repo, pageKeyService PageKeyService, glicko2 Glicko2Service) MatchHistoryServiceV2 {
	return &matchHistoryServiceV2Impl{repo, pageKeyService, glicko2}
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

// ---- proto <-> model conversions ----

func roundFromModel(r *models.Round) *pb.Round {
	ms := r.Length.Milliseconds()
	return &pb.Round{IsPlayer1Killer: &r.IsPlayer1Killer, TimeMillis: &ms}
}

func resultFromModel(r models.MatchResult) pb.MatchResult {
	switch r {
	case models.Win:
		return pb.MatchResult_MATCH_RESULT_PLAYER1_WIN
	case models.Draw:
		return pb.MatchResult_MATCH_RESULT_DRAW
	case models.Loss:
		return pb.MatchResult_MATCH_RESULT_PLAYER2_WIN
	}
	return pb.MatchResult_MATCH_RESULT_UNSPECIFIED
}

func matchFromModel(m *models.AggregatedMatch) *pb.MatchData {
	rounds := make([]*pb.Round, 0, len(m.Rounds))
	for _, r := range m.Rounds {
		rounds = append(rounds, roundFromModel(r))
	}
	unixMs := m.MatchObj.End.UTC().UnixMilli()
	return &pb.MatchData{
		Player1:      &pb.ParticipantData{Id: &m.MatchObj.Player1Id, Name: &m.MatchObj.Player1Name, Rating: &m.MatchObj.Player1Rating},
		Player2:      &pb.ParticipantData{Id: &m.MatchObj.Player2Id, Name: &m.MatchObj.Player2Name, Rating: &m.MatchObj.Player2Rating},
		Rounds:       rounds,
		EndTimestamp: &unixMs,
		MatchResult:  resultFromModel(m.MatchObj.Result),
		MatchId:      &m.MatchObj.MatchId,
	}
}

func roundToModel(r *pb.Round) (*models.Round, error) {
	if r.TimeMillis == nil {
		return nil, errors.New("no time_millis")
	}
	return &models.Round{
		IsPlayer1Killer: r.GetIsPlayer1Killer(),
		Length:          time.Duration(r.GetTimeMillis()) * time.Millisecond,
	}, nil
}

func resultToModel(r pb.MatchResult) models.MatchResult {
	switch r {
	case pb.MatchResult_MATCH_RESULT_PLAYER1_WIN:
		return models.Win
	case pb.MatchResult_MATCH_RESULT_DRAW:
		return models.Draw
	case pb.MatchResult_MATCH_RESULT_PLAYER2_WIN:
		return models.Loss
	}
	return models.Loss
}

func matchToModel(proto *pb.MatchData) (*models.AggregatedMatch, error) {
	if proto.Player1 == nil || proto.Player1.Id == nil || proto.Player1.Name == nil || proto.Player1.Rating == nil {
		return nil, errors.New("invalid player1 data")
	}
	if proto.Player2 == nil || proto.Player2.Id == nil || proto.Player2.Name == nil || proto.Player2.Rating == nil {
		return nil, errors.New("invalid player2 data")
	}
	if proto.EndTimestamp == nil {
		return nil, errors.New("no end_timestamp")
	}
	if proto.MatchResult == pb.MatchResult_MATCH_RESULT_UNSPECIFIED {
		return nil, errors.New("no match_result")
	}
	if proto.MatchId == nil {
		return nil, errors.New("no match_id")
	}
	rounds := make([]*models.Round, 0, len(proto.Rounds))
	for _, r := range proto.Rounds {
		conv, err := roundToModel(r)
		if err != nil {
			return nil, fmt.Errorf("invalid round: %w", err)
		}
		rounds = append(rounds, conv)
	}
	return &models.AggregatedMatch{
		MatchObj: &models.Match{
			MatchId:       proto.GetMatchId(),
			Player1Id:     proto.Player1.GetId(),
			Player1Name:   proto.Player1.GetName(),
			Player1Rating: proto.Player1.GetRating(),
			Player2Id:     proto.Player2.GetId(),
			Player2Name:   proto.Player2.GetName(),
			Player2Rating: proto.Player2.GetRating(),
			End:           time.UnixMilli(proto.GetEndTimestamp()),
			Result:        resultToModel(proto.MatchResult),
		},
		Rounds: rounds,
	}, nil
}

func matchResultToGameResult(r pb.MatchResult) models.GameResult {
	switch r {
	case pb.MatchResult_MATCH_RESULT_PLAYER1_WIN:
		return models.Winner
	case pb.MatchResult_MATCH_RESULT_DRAW:
		return models.DrawResult
	default:
		return models.Loser
	}
}

func ratingDataFromModel(r *models.RatingInfo) *pb.RatingData {
	if r == nil {
		return nil
	}
	display := int64(math.Floor(r.Value))
	return &pb.RatingData{RatingValue: &r.Value, DisplayValue: &display}
}

// ---- handlers ----

func (s *matchHistoryServiceV2Impl) GetMatchHistory(ctx context.Context, req *pb.GetMatchHistoryRequest) (*pb.GetMatchHistoryResponse, error) {
	logger := logging.GetLogger(ctx)
	if req.UserId == nil {
		return nil, status.Error(codes.InvalidArgument, "no user_id specified")
	}
	var matches []*models.AggregatedMatch
	var pgErr *pglib.PgLibError
	if req.Pagekey == nil {
		matches, pgErr = s.repo.GetMatchesFirstPage(ctx, req.GetUserId())
	} else {
		pageKey, err := s.pageKeyService.Deserialize(req.GetPagekey())
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid page key")
		}
		matches, pgErr = s.repo.GetMatchesSecondPage(ctx, req.GetUserId(), pageKey)
	}
	if pgErr != nil {
		logger.Error("failed to get match history", "user_id", req.GetUserId(), "error", pgErr)
		return nil, convertError(pgErr)
	}
	matchesProto := make([]*pb.MatchData, 0, len(matches))
	for _, m := range matches {
		matchesProto = append(matchesProto, matchFromModel(m))
	}
	return &pb.GetMatchHistoryResponse{
		Matches: matchesProto,
		Pagekey: s.pageKeyService.GetNewPageKey(matches),
	}, nil
}

func (s *matchHistoryServiceV2Impl) GetRating(ctx context.Context, req *pb.GetRatingRequest) (*pb.GetRatingResponse, error) {
	if req.UserId == nil {
		return nil, status.Error(codes.InvalidArgument, "no user_id specified")
	}
	initial := s.glicko2.GetInitialRating(req.GetUserId())
	rating, pgErr := s.repo.GetUserRating(ctx, initial)
	if pgErr != nil {
		return nil, convertError(pgErr)
	}
	return &pb.GetRatingResponse{Rating: ratingDataFromModel(rating)}, nil
}

func (s *matchHistoryServiceV2Impl) SaveMatch(ctx context.Context, req *pb.SaveMatchRequest) (*pb.SaveMatchResponse, error) {
	logger := logging.GetLogger(ctx)
	if req.Match == nil {
		return nil, status.Error(codes.InvalidArgument, "no match")
	}
	match, err := matchToModel(req.Match)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid match: %v", err)
	}

	initial1 := s.glicko2.GetInitialRating(match.MatchObj.Player1Id)
	rating1, pgErr := s.repo.GetUserRating(ctx, initial1)
	if pgErr != nil {
		logger.Error("failed to get player1 rating", "error", pgErr)
		return nil, convertError(pgErr)
	}
	initial2 := s.glicko2.GetInitialRating(match.MatchObj.Player2Id)
	rating2, pgErr := s.repo.GetUserRating(ctx, initial2)
	if pgErr != nil {
		logger.Error("failed to get player2 rating", "error", pgErr)
		return nil, convertError(pgErr)
	}

	gameResult := matchResultToGameResult(req.Match.MatchResult)
	newRating1, newRating2 := s.glicko2.UpdateRatings(rating1, rating2, gameResult)

	pgErr = s.repo.SaveMatch(ctx, match, []*models.RatingInfo{newRating1, newRating2})
	if pgErr != nil {
		logger.Error("failed to save match", "match_id", match.MatchObj.MatchId, "error", pgErr)
		return nil, convertError(pgErr)
	}

	return &pb.SaveMatchResponse{
		Player1Rating: ratingDataFromModel(newRating1),
		Player2Rating: ratingDataFromModel(newRating2),
	}, nil
}
