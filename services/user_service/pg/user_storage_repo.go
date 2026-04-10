package pg

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"team_dynamics/logging"
	pglib "team_dynamics/pg_lib/include"
	"team_dynamics/user_service/models"
	"time"
)

type UserStorageRepo interface {
	GetSelfData(ctx context.Context, steamId int64) (*models.UserData, *pglib.PgLibError)
	GetUserData(ctx context.Context, id int64) (*models.UserData, *pglib.PgLibError)
}

type userStorageRepoImpl struct {
	pool *pgxpool.Pool
}

func MakeUserStorageRepo(pool *pgxpool.Pool) UserStorageRepo {
	return userStorageRepoImpl{pool}
}

func getSelfDataImpl(ctx context.Context, tx pgx.Tx, steamId int64) (*models.UserData, error, pglib.ResponseStatus) {
	logger := logging.GetLogger(ctx)
	var res models.UserData
	err := tx.QueryRow(ctx, `
SELECT id, name FROM users
WHERE steam_id = $1
`, steamId).Scan(&res.Id, &res.Name)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			logger.Debug("cannot get user id", "error", err)
			return nil, err, pglib.NormalRetry
		}
	}
	if err == nil {
		return &res, nil, pglib.NoRetry
	}

	res.Id = steamId
	res.Name = uuid.New().String()

	_, err = tx.Exec(ctx, `
INSERT INTO users (id, name, steam_id) VALUES ($1, $2, $3)
`, res.Id, res.Name, steamId)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pglib.IsSerializationError(pgErr) {
				logger.Debug("Serialization error while inserting user", "error", err)
				return nil, err, pglib.FreeRetry
			}
		}
		logger.Debug("unknown error while inserting user", "error", err)
		return nil, err, pglib.NormalRetry
	}
	return &res, nil, pglib.NoRetry
}

func getUserData(ctx context.Context, tx pgx.Tx, id int64) (*models.UserData, error, pglib.ResponseStatus) {
	logger := logging.GetLogger(ctx)
	var res models.UserData
	err := tx.QueryRow(ctx, `
SELECT id, name FROM users
WHERE id = $1
`, id).Scan(&res.Id, &res.Name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Info("No user found")
			return nil, err, pglib.NoRetry
		}
		logger.Debug("unknown error while getting user", "error", err)
		return nil, err, pglib.NormalRetry
	}
	return &res, nil, pglib.NoRetry
}

func (r userStorageRepoImpl) GetSelfData(ctx context.Context, steamId int64) (*models.UserData, *pglib.PgLibError) {
	return pglib.PerformOperation(ctx, r.pool, &pglib.QueryConfig{
		Retries:        3,
		Timeout:        200 * time.Millisecond,
		IsolationLevel: pgx.RepeatableRead,
		AccessMode:     pgx.ReadWrite,
	}, func(ctx1 context.Context, tx pgx.Tx) (*models.UserData, error, pglib.ResponseStatus) {
		return getSelfDataImpl(ctx1, tx, steamId)
	})
}

func (r userStorageRepoImpl) GetUserData(ctx context.Context, id int64) (*models.UserData, *pglib.PgLibError) {
	return pglib.PerformOperation(ctx, r.pool, &pglib.QueryConfig{
		Retries:        3,
		Timeout:        200 * time.Millisecond,
		IsolationLevel: pgx.ReadCommitted,
		AccessMode:     pgx.ReadOnly,
	}, func(ctx1 context.Context, tx pgx.Tx) (*models.UserData, error, pglib.ResponseStatus) {
		return getUserData(ctx1, tx, id)
	})
}
