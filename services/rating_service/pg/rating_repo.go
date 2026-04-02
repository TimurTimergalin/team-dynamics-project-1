package pg

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	pglib "team_dynamics/pg_lib/include"
	"team_dynamics/rating_service/models"
	"time"
)

type RatingRepo interface {
	GetUserRating(ctx context.Context, ratingInfo *models.RatingInfo) (*models.RatingInfo, *pglib.PgLibError)
	UpdateUserRating(ctx context.Context, ratingInfo []*models.RatingInfo, matchId string) *pglib.PgLibError
}

type ratingRepoImpl struct {
	pool *pgxpool.Pool
}

func getUserRatingImpl(ctx context.Context, pool *pgxpool.Pool, ratingInfo *models.RatingInfo) (*models.RatingInfo, error, pglib.RetryPolicy) {
	var ratingInfoRes models.RatingInfo
	err := pool.QueryRow(ctx, "SELECT id, rating, rating_distribution FROM ratings WHERE id = $1", ratingInfo.UserId).Scan(&ratingInfoRes.UserId, &ratingInfoRes.Value, &ratingInfoRes.Deviation)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			_, err := pool.Exec(ctx, "INSERT INTO ratings (id, rating, rating_distribution) VALUES ($1, $2, $3)", ratingInfo.UserId, ratingInfo.Value, ratingInfo.Deviation)
			if err != nil {
				var pgErr *pgconn.PgError
				if errors.As(err, &pgErr) {
					return nil, err, pglib.FreeRetry
				}
				return nil, err, pglib.NormalRetry
			}
			return ratingInfo, nil, pglib.NoRetry
		}
		return nil, err, pglib.NormalRetry
	}
	return &ratingInfoRes, nil, pglib.NoRetry
}

func updateUserRatingImpl(ctx context.Context, pool *pgxpool.Pool, ratingInfos []*models.RatingInfo, matchId string) (*struct{}, error, pglib.RetryPolicy) {
	_, err := pool.Exec(ctx, "INSERT INTO matches (matchId) VALUES ($1)", matchId)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pglib.IsConstraintViolated(pgErr) {
				return nil, nil, pglib.NoRetry
			}
			if pglib.IsSerializationError(pgErr) {
				return nil, err, pglib.FreeRetry
			}
		}
		return nil, err, pglib.NormalRetry
	}
	for _, ratingInfo := range ratingInfos {
		var foundUserId int64
		err := pool.QueryRow(ctx, "UPDATE ratings SET (rating, rating_distribution) = ($2, $3) WHERE id = $1 RETURNING id", ratingInfo.UserId, ratingInfo.Value, ratingInfo.Deviation).Scan(&foundUserId)
		if err != nil {
			if pglib.IsNoRows(err) {
				return nil, err, pglib.NoRetry
			}
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				if pglib.IsSerializationError(pgErr) {
					return nil, err, pglib.FreeRetry
				}
			}
			return nil, err, pglib.NormalRetry
		}
	}
	return nil, nil, pglib.NoRetry
}

func (repo ratingRepoImpl) GetUserRating(ctx context.Context, ratingInfo *models.RatingInfo) (*models.RatingInfo, *pglib.PgLibError) {
	return pglib.PerformOperation(ctx, repo.pool, &pglib.QueryConfig{
		Retries:        3,
		Timeout:        50 * time.Millisecond,
		IsolationLevel: pgx.RepeatableRead,
		AccessMode:     pgx.ReadWrite,
	}, func(ctx1 context.Context, pool1 *pgxpool.Pool) (*models.RatingInfo, error, pglib.RetryPolicy) {
		return getUserRatingImpl(ctx1, pool1, ratingInfo)
	})
}

func (repo ratingRepoImpl) UpdateUserRating(ctx context.Context, updateRequests []*models.RatingInfo, matchId string) *pglib.PgLibError {
	_, err := pglib.PerformOperation(ctx, repo.pool, &pglib.QueryConfig{
		Retries:        3,
		Timeout:        100 * time.Millisecond,
		IsolationLevel: pgx.RepeatableRead,
		AccessMode:     pgx.ReadWrite,
	}, func(ctx1 context.Context, pool1 *pgxpool.Pool) (*struct{}, error, pglib.RetryPolicy) {
		return updateUserRatingImpl(ctx1, pool1, updateRequests, matchId)
	})
	return err
}

func MakeRatingRepo(pool *pgxpool.Pool) RatingRepo {
	return ratingRepoImpl{pool: pool}
}
