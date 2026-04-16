package pg

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"team_dynamics/logging"
	"team_dynamics/match_history_service/models"
	pglib "team_dynamics/pg_lib/include"
	"time"
)

type MatchHistoryRepo interface {
	GetMatchesFirstPage(ctx context.Context, userId int64) ([]*models.AggregatedMatch, *pglib.PgLibError)
	GetMatchesSecondPage(ctx context.Context, userId int64, pageKey *models.PageKey) ([]*models.AggregatedMatch, *pglib.PgLibError)
	SaveMatch(ctx context.Context, match *models.AggregatedMatch) *pglib.PgLibError
}

type matchHistoryRepoImpl struct {
	pool *pgxpool.Pool
}

func getAggregatedMatches(ctx context.Context, tx pgx.Tx, matchRows pgx.Rows) (*[]*models.AggregatedMatch, error, pglib.ResponseStatus) {
	logger := logging.GetLogger(ctx)
	matches := make([]*models.Match, 0, 10)
	for matchRows.Next() {
		var m models.Match
		err := matchRows.Scan(&m.MatchId, &m.Player1Id, &m.Player2Id, &m.Player1Name, &m.Player2Name, &m.Player1Rating, &m.Player2Rating, &m.End, &m.Result)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				if pglib.IsServerSideTimeout(pgErr) {
					logger.Debug("Timeout when getting match list (server)", "error", err)
					return nil, err, pglib.NormalRetry
				}
			}
			if errors.Is(err, context.DeadlineExceeded) {
				logger.Debug("Timeout when getting match list (client)", "error", err)
				return nil, err, pglib.NormalRetry
			}
			logger.Error("Error occurred wile reading match", "error", err)
			return nil, err, pglib.NoRetry
		}
		matches = append(matches, &m)
	}
	if len(matches) == 0 {
		result := make([]*models.AggregatedMatch, 0)
		return &result, nil, pglib.NoRetry
	}
	matchIds := make([]string, 0, len(matches))
	for _, match := range matches {
		matchIds = append(matchIds, match.MatchId)
	}
	roundRows, err := tx.Query(ctx, `
SELECT match_id, is_player1_killer, time_millis
FROM rounds
WHERE match_id = ANY($1)
ORDER BY round_number ASC
`, matchIds)
	if err != nil {
		logger.Debug("Cannot get round list", "error", err)
		return nil, err, pglib.NormalRetry
	}
	defer roundRows.Close()
	roundsMap := make(map[string][]*models.Round)
	for roundRows.Next() {
		var r models.Round
		var matchId string
		err := roundRows.Scan(&matchId, &r.IsPlayer1Killer, &r.Length)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				if pglib.IsServerSideTimeout(pgErr) {
					logger.Debug("Timeout when getting round list (server)", "error", err)
					return nil, err, pglib.NormalRetry
				}
			}
			if errors.Is(err, context.DeadlineExceeded) {
				logger.Debug("Timeout when getting round list (client)", "error", err)
				return nil, err, pglib.NormalRetry
			}
			logger.Error("Error occurred wile round match", "error", err)
			return nil, err, pglib.NoRetry
		}
		roundsMap[matchId] = append(roundsMap[matchId], &r)
	}
	results := make([]*models.AggregatedMatch, 0, len(matches))
	for _, match := range matches {
		matchId := match.MatchId
		rounds, ok := roundsMap[matchId]
		if !ok {
			rounds = make([]*models.Round, 0)
		}
		results = append(results, &models.AggregatedMatch{
			MatchObj: match,
			Rounds:   rounds,
		})
	}
	return &results, nil, pglib.NoRetry
}

func getMatchesFirstPageImpl(ctx context.Context, tx pgx.Tx, userId int64) (*[]*models.AggregatedMatch, error, pglib.ResponseStatus) {
	logger := logging.GetLogger(ctx)
	matchRows, err := tx.Query(ctx, `
SELECT match_id, player1_id, player2_id, player1_name, player2_name,
player1_rating, player2_rating, end_timestamp, result
FROM matches
WHERE player1_id = $1 OR player2_id = $1
ORDER BY end_timestamp DESC
LIMIT 10
`, userId)
	if err != nil {
		logger.Debug("Cannot get match list", "error", err)
		return nil, err, pglib.NormalRetry
	}
	defer matchRows.Close()
	return getAggregatedMatches(ctx, tx, matchRows)
}

func getMatchesSecondPageImpl(ctx context.Context, tx pgx.Tx, userId int64, pageKey *models.PageKey) (*[]*models.AggregatedMatch, error, pglib.ResponseStatus) {
	logger := logging.GetLogger(ctx)
	matchRows, err := tx.Query(ctx, `
SELECT match_id, player1_id, player2_id, player1_name, player2_name,
player1_rating, player2_rating, end_timestamp, result
FROM matches
WHERE (player1_id = $1 OR player2_id = $1) AND end_timestamp < $2
ORDER BY end_timestamp DESC
LIMIT 10
`, userId, pageKey.Before)
	if err != nil {
		logger.Debug("Cannot get match list", "error", err)
		return nil, err, pglib.NormalRetry
	}
	defer matchRows.Close()
	return getAggregatedMatches(ctx, tx, matchRows)
}

func saveMatchImpl(ctx context.Context, tx pgx.Tx, match *models.AggregatedMatch) (*struct{}, error, pglib.ResponseStatus) {
	matchObj := match.MatchObj
	logger := logging.GetLogger(ctx)
	_, err := tx.Exec(ctx, `
INSERT INTO matches (match_id, player1_id, player2_id, player1_name, player2_name, player1_rating, player2_rating, end_timestamp, result)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
`, matchObj.MatchId, matchObj.Player1Id, matchObj.Player2Id, matchObj.Player1Name, matchObj.Player2Name, matchObj.Player1Rating, matchObj.Player2Rating, matchObj.End, matchObj.Result)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pglib.IsConstraintViolated(pgErr) {
				logger.Info("Constraint violated when inserting match", "error", err)
				return nil, nil, pglib.ForceRollback
			}
			if pglib.IsSerializationError(pgErr) {
				logger.Info("Serialization error encountered when inserting match", "error", err)
				return nil, err, pglib.FreeRetry
			}
		}
		logger.Debug("Retrying inserting match")
		return nil, err, pglib.NormalRetry
	}
	rounds := match.Rounds
	if len(rounds) > 0 {
		batch := &pgx.Batch{}
		for i, r := range rounds {
			batch.Queue(`
INSERT INTO rounds (match_id, round_number, is_player1_killer, time_millis)
VALUES ($1, $2, $3, $4)
`, matchObj.MatchId, i, r.IsPlayer1Killer, r.Length.Milliseconds())
		}
		br := tx.SendBatch(ctx, batch)
		defer func(br pgx.BatchResults) {
			_ = br.Close()
		}(br)

		for i := 0; i < len(rounds); i += 1 {
			_, err := br.Exec()
			if err != nil {
				var pgErr *pgconn.PgError
				if errors.As(err, &pgErr) {
					if pglib.IsConstraintViolated(pgErr) {
						logger.Info("Constraint violated when inserting round", "error", err)
						return nil, err, pglib.NoRetry
					}
					if pglib.IsSerializationError(pgErr) {
						logger.Info("Serialization error encountered when inserting round", "error", err)
						return nil, err, pglib.FreeRetry
					}
				}
				logger.Debug("Retrying inserting round")
				return nil, err, pglib.NormalRetry
			}
		}
	}

	return nil, nil, pglib.NoRetry
}

func (r *matchHistoryRepoImpl) GetMatchesFirstPage(ctx context.Context, userId int64) ([]*models.AggregatedMatch, *pglib.PgLibError) {
	matches, err := pglib.PerformOperation(ctx, r.pool, &pglib.QueryConfig{
		Retries:        3,
		Timeout:        200 * time.Millisecond,
		IsolationLevel: pgx.ReadCommitted,
		AccessMode:     pgx.ReadOnly,
	}, func(ctx1 context.Context, tx pgx.Tx) (*[]*models.AggregatedMatch, error, pglib.ResponseStatus) {
		return getMatchesFirstPageImpl(ctx1, tx, userId)
	})
	if matches == nil {
		return make([]*models.AggregatedMatch, 0), err
	}
	return *matches, err
}

func (r *matchHistoryRepoImpl) GetMatchesSecondPage(ctx context.Context, userId int64, pageKey *models.PageKey) ([]*models.AggregatedMatch, *pglib.PgLibError) {
	matches, err := pglib.PerformOperation(ctx, r.pool, &pglib.QueryConfig{
		Retries:        3,
		Timeout:        200 * time.Millisecond,
		IsolationLevel: pgx.ReadCommitted,
		AccessMode:     pgx.ReadOnly,
	}, func(ctx1 context.Context, tx pgx.Tx) (*[]*models.AggregatedMatch, error, pglib.ResponseStatus) {
		return getMatchesSecondPageImpl(ctx1, tx, userId, pageKey)
	})
	if matches == nil {
		return make([]*models.AggregatedMatch, 0), err
	}
	return *matches, err
}

func (r *matchHistoryRepoImpl) SaveMatch(ctx context.Context, match *models.AggregatedMatch) *pglib.PgLibError {
	_, err := pglib.PerformOperation(ctx, r.pool, &pglib.QueryConfig{
		Retries:        3,
		Timeout:        200 * time.Millisecond,
		IsolationLevel: pgx.RepeatableRead,
		AccessMode:     pgx.ReadWrite,
	}, func(ctx1 context.Context, tx pgx.Tx) (*struct{}, error, pglib.ResponseStatus) {
		return saveMatchImpl(ctx1, tx, match)
	})
	return err
}

func MakeMatchHistoryRepo(pool *pgxpool.Pool) MatchHistoryRepo {
	return &matchHistoryRepoImpl{pool}
}
