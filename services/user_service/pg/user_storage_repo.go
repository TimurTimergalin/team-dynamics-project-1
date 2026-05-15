package pg

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"team_dynamics/logging"
	pglib "team_dynamics/pg_lib/include"
	"team_dynamics/user_service/models"
	"time"
)

type UserStorageRepo interface {
	UpsertSelfData(ctx context.Context, steamId int64, name string) (*models.UserData, *pglib.PgLibError)
	UpsertSelfDataEos(ctx context.Context, eosId string, name string) (*models.UserData, *pglib.PgLibError)
	GetUserData(ctx context.Context, id int64) (*models.UserData, *pglib.PgLibError)
	GetFriends(ctx context.Context, userId int64, key *models.PageKey) ([]*models.Friend, *pglib.PgLibError)
	GetIncomingRequests(ctx context.Context, userId int64, key *models.PageKey) ([]*models.Friend, *pglib.PgLibError)
	GetOutgoingRequests(ctx context.Context, userId int64, key *models.PageKey) ([]*models.Friend, *pglib.PgLibError)
	AddFriend(ctx context.Context, userId int64, otherUserId int64) (*models.Friend, models.AddFriendResult, *pglib.PgLibError)
	RemoveFriend(ctx context.Context, userId int64, otherUserId int64) (models.RemoveFriendResult, *pglib.PgLibError)
}

type userStorageRepoImpl struct {
	pool *pgxpool.Pool
}

func MakeUserStorageRepo(pool *pgxpool.Pool) UserStorageRepo {
	return userStorageRepoImpl{pool}
}

func upsertSelfDataImpl(ctx context.Context, tx pgx.Tx, steamId int64, name string) (*models.UserData, error, pglib.ResponseStatus) {
	logger := logging.GetLogger(ctx)
	var res models.UserData
	err := tx.QueryRow(ctx, `
INSERT INTO users (name, steam_id) VALUES ($1, $2)
ON CONFLICT (steam_id) DO UPDATE SET name = EXCLUDED.name
RETURNING id, name
`, name, steamId).Scan(&res.Id, &res.Name)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pglib.IsSerializationError(pgErr) {
				logger.Debug("serialization error while upserting user", "error", err)
				return nil, err, pglib.FreeRetry
			}
		}
		logger.Debug("error while upserting user", "error", err)
		return nil, err, pglib.NormalRetry
	}
	return &res, nil, pglib.NoRetry
}

func upsertSelfDataEosImpl(ctx context.Context, tx pgx.Tx, eosId string, name string) (*models.UserData, error, pglib.ResponseStatus) {
	logger := logging.GetLogger(ctx)
	var res models.UserData
	err := tx.QueryRow(ctx, `
INSERT INTO users (name, eos_id) VALUES ($1, $2)
ON CONFLICT (eos_id) DO UPDATE SET name = EXCLUDED.name
RETURNING id, name
`, name, eosId).Scan(&res.Id, &res.Name)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pglib.IsSerializationError(pgErr) {
				logger.Debug("serialization error while upserting eos user", "error", err)
				return nil, err, pglib.FreeRetry
			}
		}
		logger.Debug("error while upserting eos user", "error", err)
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

func readFriendList(ctx context.Context, rows pgx.Rows) ([]*models.Friend, error, pglib.ResponseStatus) {
	logger := logging.GetLogger(ctx)
	result := make([]*models.Friend, 0, 20)
	for rows.Next() {
		var friend models.Friend
		friend.Data = &models.UserData{}
		err := rows.Scan(&friend.Data.Id, &friend.Data.Name, &friend.Source)
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
		result = append(result, &friend)
	}
	return result, nil, pglib.NoRetry
}

func getFriendFirstPage(ctx context.Context, tx pgx.Tx, userId int64) ([]*models.Friend, error, pglib.ResponseStatus) {
	logger := logging.GetLogger(ctx)
	rows, err := tx.Query(ctx, `
SELECT u.id, u.name, f.source
FROM (
    SELECT player2_id AS friend_id, source
    FROM friends
    WHERE player1_id = $1
    
    UNION ALL
    
    SELECT player1_id AS friend_id, source
    FROM friends
    WHERE player2_id = $1
) AS f
JOIN users u ON u.id = f.friend_id
ORDER BY u.id ASC
LIMIT 20;
`, userId)
	if err != nil {
		logger.Debug("Cannot get friend list", "error", err)
		return nil, err, pglib.NormalRetry
	}
	defer rows.Close()
	return readFriendList(ctx, rows)
}

func getFriendsSecondPage(ctx context.Context, tx pgx.Tx, userId int64, pageKey *models.PageKey) ([]*models.Friend, error, pglib.ResponseStatus) {
	logger := logging.GetLogger(ctx)
	rows, err := tx.Query(ctx, `
SELECT u.id, u.name, f.source
FROM (
    SELECT player2_id AS friend_id, source
    FROM friends
    WHERE player1_id = $1
      AND player2_id > $2
    
    UNION ALL
    
    SELECT player1_id AS friend_id, source
    FROM friends
    WHERE player2_id = $1
      AND player1_id > $2
) AS f
JOIN users u ON u.id = f.friend_id
ORDER BY u.id ASC
LIMIT 20;
`, userId, pageKey.LastUserId)
	if err != nil {
		logger.Debug("Cannot get friend list", "error", err)
		return nil, err, pglib.NormalRetry
	}
	defer rows.Close()
	return readFriendList(ctx, rows)
}

func getFriendsImpl(ctx context.Context, tx pgx.Tx, userId int64, pageKey *models.PageKey) ([]*models.Friend, error, pglib.ResponseStatus) {
	if pageKey != nil {
		return getFriendsSecondPage(ctx, tx, userId, pageKey)
	}
	return getFriendFirstPage(ctx, tx, userId)
}

func getIncomingRequests(ctx context.Context, tx pgx.Tx, userId int64, pageKey *models.PageKey) ([]*models.Friend, error, pglib.ResponseStatus) {
	logger := logging.GetLogger(ctx)
	rows, err := tx.Query(ctx, `
SELECT u.id, u.name, r.source
FROM requests r
JOIN users u ON u.id = r.from
WHERE r.to = $1
  AND r.from > $2
ORDER BY r.from ASC
LIMIT 20;
`, userId, pageKey.LastUserId)
	if err != nil {
		logger.Debug("Cannot get incoming requests list", "error", err)
		return nil, err, pglib.NormalRetry
	}
	defer rows.Close()
	return readFriendList(ctx, rows)
}

func getOutgoingRequests(ctx context.Context, tx pgx.Tx, userId int64, pageKey *models.PageKey) ([]*models.Friend, error, pglib.ResponseStatus) {
	logger := logging.GetLogger(ctx)
	rows, err := tx.Query(ctx, `
SELECT u.id, u.name, r.source
FROM requests r
JOIN users u ON u.id = r.to
WHERE r.from = $1
  AND r.to > $2
ORDER BY r.to ASC
LIMIT 20;
`, userId, pageKey.LastUserId)
	if err != nil {
		logger.Debug("Cannot get outgoing requests list", "error", err)
		return nil, err, pglib.NormalRetry
	}
	defer rows.Close()
	return readFriendList(ctx, rows)
}

func (r userStorageRepoImpl) UpsertSelfData(ctx context.Context, steamId int64, name string) (*models.UserData, *pglib.PgLibError) {
	return pglib.PerformOperation(ctx, r.pool, &pglib.QueryConfig{
		Retries:        3,
		Timeout:        200 * time.Millisecond,
		IsolationLevel: pgx.RepeatableRead,
		AccessMode:     pgx.ReadWrite,
	}, func(ctx1 context.Context, tx pgx.Tx) (*models.UserData, error, pglib.ResponseStatus) {
		return upsertSelfDataImpl(ctx1, tx, steamId, name)
	})
}

func (r userStorageRepoImpl) UpsertSelfDataEos(ctx context.Context, eosId string, name string) (*models.UserData, *pglib.PgLibError) {
	return pglib.PerformOperation(ctx, r.pool, &pglib.QueryConfig{
		Retries:        3,
		Timeout:        200 * time.Millisecond,
		IsolationLevel: pgx.RepeatableRead,
		AccessMode:     pgx.ReadWrite,
	}, func(ctx1 context.Context, tx pgx.Tx) (*models.UserData, error, pglib.ResponseStatus) {
		return upsertSelfDataEosImpl(ctx1, tx, eosId, name)
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

func (r userStorageRepoImpl) GetFriends(ctx context.Context, userId int64, key *models.PageKey) ([]*models.Friend, *pglib.PgLibError) {
	return pglib.PerformOperation(ctx, r.pool, &pglib.QueryConfig{
		Retries:        3,
		Timeout:        200 * time.Millisecond,
		IsolationLevel: pgx.ReadCommitted,
		AccessMode:     pgx.ReadOnly,
	}, func(ctx1 context.Context, tx pgx.Tx) ([]*models.Friend, error, pglib.ResponseStatus) {
		return getFriendsImpl(ctx1, tx, userId, key)
	})
}

func (r userStorageRepoImpl) GetIncomingRequests(ctx context.Context, userId int64, key *models.PageKey) ([]*models.Friend, *pglib.PgLibError) {
	return pglib.PerformOperation(ctx, r.pool, &pglib.QueryConfig{
		Retries:        3,
		Timeout:        200 * time.Millisecond,
		IsolationLevel: pgx.ReadCommitted,
		AccessMode:     pgx.ReadOnly,
	}, func(ctx1 context.Context, tx pgx.Tx) ([]*models.Friend, error, pglib.ResponseStatus) {
		return getIncomingRequests(ctx1, tx, userId, key)
	})
}

func (r userStorageRepoImpl) GetOutgoingRequests(ctx context.Context, userId int64, key *models.PageKey) ([]*models.Friend, *pglib.PgLibError) {
	return pglib.PerformOperation(ctx, r.pool, &pglib.QueryConfig{
		Retries:        3,
		Timeout:        200 * time.Millisecond,
		IsolationLevel: pgx.ReadCommitted,
		AccessMode:     pgx.ReadOnly,
	}, func(ctx1 context.Context, tx pgx.Tx) ([]*models.Friend, error, pglib.ResponseStatus) {
		return getOutgoingRequests(ctx1, tx, userId, key)
	})
}

func addFriendImpl(ctx context.Context, tx pgx.Tx, userId int64, otherUserId int64) (*models.Friend, models.AddFriendResult, error, pglib.ResponseStatus) {
	logger := logging.GetLogger(ctx)

	// Verify both users exist
	var exists bool
	err := tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`, userId).Scan(&exists)
	if err != nil {
		return nil, models.AddFriendNoop, err, pglib.NormalRetry
	}
	if !exists {
		return nil, models.AddFriendNoop, fmt.Errorf("user %d not found", userId), pglib.NoRetry
	}
	err = tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`, otherUserId).Scan(&exists)
	if err != nil {
		return nil, models.AddFriendNoop, err, pglib.NormalRetry
	}
	if !exists {
		return nil, models.AddFriendNoop, fmt.Errorf("user %d not found", otherUserId), pglib.NoRetry
	}

	// Check requests table
	var fromSource int64
	var hasOutgoing, hasIncoming bool
	err = tx.QueryRow(ctx, `SELECT source FROM requests WHERE "from" = $1 AND "to" = $2`, userId, otherUserId).Scan(&fromSource)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, models.AddFriendNoop, err, pglib.NormalRetry
	}
	hasOutgoing = err == nil

	var incomingSource int64
	err = tx.QueryRow(ctx, `SELECT source FROM requests WHERE "from" = $1 AND "to" = $2`, otherUserId, userId).Scan(&incomingSource)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, models.AddFriendNoop, err, pglib.NormalRetry
	}
	hasIncoming = err == nil

	if hasOutgoing {
		return nil, models.AddFriendNoop, nil, pglib.NoRetry
	}

	// Check if already friends
	p1, p2 := userId, otherUserId
	if p1 > p2 {
		p1, p2 = p2, p1
	}
	var alreadyFriends bool
	err = tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM friends WHERE player1_id = $1 AND player2_id = $2)`, p1, p2).Scan(&alreadyFriends)
	if err != nil {
		return nil, models.AddFriendNoop, err, pglib.NormalRetry
	}
	if alreadyFriends {
		return nil, models.AddFriendNoop, nil, pglib.NoRetry
	}

	if hasIncoming {
		_, err = tx.Exec(ctx, `
INSERT INTO friends (player1_id, player2_id, source)
VALUES ($1, $2, $3)
ON CONFLICT (player1_id, player2_id) DO NOTHING
`, p1, p2, incomingSource)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				if pglib.IsSerializationError(pgErr) {
					return nil, models.AddFriendNoop, err, pglib.FreeRetry
				}
			}
			logger.Debug("failed to insert friend", "error", err)
			return nil, models.AddFriendNoop, err, pglib.NormalRetry
		}
		_, err = tx.Exec(ctx, `DELETE FROM requests WHERE "from" = $1 AND "to" = $2`, otherUserId, userId)
		if err != nil {
			logger.Debug("failed to delete request", "error", err)
			return nil, models.AddFriendNoop, err, pglib.NormalRetry
		}
		var otherUserData models.UserData
		err = tx.QueryRow(ctx, `SELECT id, name FROM users WHERE id = $1`, otherUserId).Scan(&otherUserData.Id, &otherUserData.Name)
		if err != nil {
			return nil, models.AddFriendNoop, err, pglib.NormalRetry
		}
		return &models.Friend{Data: &otherUserData, Source: models.FriendSource(incomingSource)}, models.AddFriendAccepted, nil, pglib.NoRetry
	}

	_, err = tx.Exec(ctx, `INSERT INTO requests ("from", "to", source) VALUES ($1, $2, $3)`, userId, otherUserId, int64(models.Internal))
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pglib.IsSerializationError(pgErr) {
				return nil, models.AddFriendNoop, err, pglib.FreeRetry
			}
		}
		logger.Debug("failed to insert request", "error", err)
		return nil, models.AddFriendNoop, err, pglib.NormalRetry
	}
	return nil, models.AddFriendRequestSent, nil, pglib.NoRetry
}

func removeFriendImpl(ctx context.Context, tx pgx.Tx, userId int64, otherUserId int64) (models.RemoveFriendResult, error, pglib.ResponseStatus) {
	logger := logging.GetLogger(ctx)

	p1, p2 := userId, otherUserId
	if p1 > p2 {
		p1, p2 = p2, p1
	}
	tag, err := tx.Exec(ctx, `DELETE FROM friends WHERE player1_id = $1 AND player2_id = $2`, p1, p2)
	if err != nil {
		logger.Debug("failed to delete friend", "error", err)
		return models.RemoveFriendNoop, err, pglib.NormalRetry
	}
	if tag.RowsAffected() > 0 {
		return models.RemoveFriendFriendRemoved, nil, pglib.NoRetry
	}

	// No friend row — check for outgoing request (userId → otherUserId)
	tag, err = tx.Exec(ctx, `DELETE FROM requests WHERE "from" = $1 AND "to" = $2`, userId, otherUserId)
	if err != nil {
		logger.Debug("failed to delete outgoing request", "error", err)
		return models.RemoveFriendNoop, err, pglib.NormalRetry
	}
	if tag.RowsAffected() > 0 {
		return models.RemoveFriendRequestCancelled, nil, pglib.NoRetry
	}

	// Check for incoming request (otherUserId → userId)
	tag, err = tx.Exec(ctx, `DELETE FROM requests WHERE "from" = $1 AND "to" = $2`, otherUserId, userId)
	if err != nil {
		logger.Debug("failed to delete incoming request", "error", err)
		return models.RemoveFriendNoop, err, pglib.NormalRetry
	}
	if tag.RowsAffected() > 0 {
		return models.RemoveFriendRequestDeclined, nil, pglib.NoRetry
	}

	return models.RemoveFriendNoop, nil, pglib.NoRetry
}

type addFriendResult struct {
	friend *models.Friend
	result models.AddFriendResult
}

func (r userStorageRepoImpl) AddFriend(ctx context.Context, userId int64, otherUserId int64) (*models.Friend, models.AddFriendResult, *pglib.PgLibError) {
	res, pgErr := pglib.PerformOperation(ctx, r.pool, &pglib.QueryConfig{
		Retries:        3,
		Timeout:        200 * time.Millisecond,
		IsolationLevel: pgx.RepeatableRead,
		AccessMode:     pgx.ReadWrite,
	}, func(ctx1 context.Context, tx pgx.Tx) (addFriendResult, error, pglib.ResponseStatus) {
		friend, result, err, status := addFriendImpl(ctx1, tx, userId, otherUserId)
		return addFriendResult{friend, result}, err, status
	})
	return res.friend, res.result, pgErr
}

func (r userStorageRepoImpl) RemoveFriend(ctx context.Context, userId int64, otherUserId int64) (models.RemoveFriendResult, *pglib.PgLibError) {
	return pglib.PerformOperation(ctx, r.pool, &pglib.QueryConfig{
		Retries:        3,
		Timeout:        200 * time.Millisecond,
		IsolationLevel: pgx.RepeatableRead,
		AccessMode:     pgx.ReadWrite,
	}, func(ctx1 context.Context, tx pgx.Tx) (models.RemoveFriendResult, error, pglib.ResponseStatus) {
		return removeFriendImpl(ctx1, tx, userId, otherUserId)
	})
}
