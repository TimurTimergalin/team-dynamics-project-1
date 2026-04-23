package pg_lib

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"team_dynamics/logging"
)

type ResponseStatus = int32

const (
	NoRetry ResponseStatus = iota
	NormalRetry
	FreeRetry
	ForceRollback
)

func withTransaction[T any](ctx context.Context, pool *pgxpool.Pool, cfg *QueryConfig, op func(context.Context, pgx.Tx) (_res T, _err error, _shouldRetry ResponseStatus)) (_res T, _err *PgLibError, _shouldRetry ResponseStatus) {
	var zero T
	logger := logging.GetLogger(ctx)
	tx, err := pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   cfg.IsolationLevel,
		AccessMode: cfg.AccessMode,
	})
	if err != nil {
		logger.Debug("Could not begin tx")
		return zero, makeConnectionError(err), NormalRetry
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		_ = tx.Rollback(ctx)
	}(tx, ctx)

	res, err, shouldRetry := op(ctx, tx)
	if err != nil {
		if shouldRetry == NoRetry {
			logger.Debug("Some non-retryable error occurred", "error", err)
			return zero, makeLogicError(err), shouldRetry
		}
		logger.Debug("Some retryable error occurred", "error", err)
		return zero, makeConnectionError(err), shouldRetry
	}
	if shouldRetry == ForceRollback {
		logger.Debug("Rollback without error")
		return zero, nil, NoRetry
	}

	err = tx.Commit(ctx)
	if err != nil {
		logger.Debug("Could not commit", "error", err)
		return zero, makeConnectionError(err), NormalRetry
	}

	return res, nil, NoRetry
}

func withTimeout[T any](ctx context.Context, cfg *QueryConfig, op func(context.Context) (_res T, _err *PgLibError, _shouldRetry ResponseStatus)) (_res T, _err *PgLibError, _shouldRetry ResponseStatus) {
	ctx, cancel := context.WithTimeout(ctx, cfg.Timeout)
	defer cancel()
	return op(ctx)
}

func withRetry[T any](ctx context.Context, cfg *QueryConfig, op func() (_res T, _err *PgLibError, _shouldRetry ResponseStatus)) (T, *PgLibError) {
	var zero T
	logger := logging.GetLogger(ctx)
	var lastErr *PgLibError = nil
	maxAllowedFree := 5
	for i := cfg.Retries; i > 0; i -= 1 {
		res, err, shouldRetry := op()
		if err == nil {
			return res, nil
		}
		if shouldRetry == NoRetry {
			return zero, err
		}
		if shouldRetry == FreeRetry && maxAllowedFree > 0 {
			i += 1
			maxAllowedFree -= 1
		}
		lastErr = err
	}
	if lastErr == nil {
		return zero, makeServerError(fmt.Errorf("no operation was performed"))
	}
	logger.Debug("All retries failed", "error", lastErr)
	return zero, makeConnectionError(lastErr)
}

func PerformOperation[T any](ctx context.Context, pool *pgxpool.Pool, cfg *QueryConfig, op func(context.Context, pgx.Tx) (T, error, ResponseStatus)) (T, *PgLibError) {
	res, err := withRetry[T](ctx, cfg, func() (T, *PgLibError, ResponseStatus) {
		return withTimeout[T](ctx, cfg, func(timeoutCtx context.Context) (T, *PgLibError, ResponseStatus) {
			return withTransaction[T](timeoutCtx, pool, cfg, op)
		})
	})
	return res, err
}
