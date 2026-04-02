package pg_lib

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RetryPolicy = int32

const (
	NoRetry RetryPolicy = iota
	NormalRetry
	FreeRetry
)

func withTransaction[T any](ctx context.Context, pool *pgxpool.Pool, cfg *QueryConfig, op func(context.Context, *pgxpool.Pool) (_res *T, _err error, _shouldRetry RetryPolicy)) (_res *T, _err *PgLibError, _shouldRetry RetryPolicy) {
	tx, err := pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   cfg.IsolationLevel,
		AccessMode: cfg.AccessMode,
	})
	if err != nil {
		return nil, makeConnectionError(err), NormalRetry
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		_ = tx.Rollback(ctx)
	}(tx, ctx)

	res, err, shouldRetry := op(ctx, pool)
	if err != nil {
		if shouldRetry == NoRetry {
			return nil, makeLogicError(err), shouldRetry
		}
		return nil, makeConnectionError(err), shouldRetry
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, makeConnectionError(err), NormalRetry
	}

	return res, nil, NoRetry
}

func withTimeout[T any](ctx context.Context, cfg *QueryConfig, op func(context.Context) (_res *T, _err *PgLibError, _shouldRetry RetryPolicy)) (_res *T, _err *PgLibError, _shouldRetry RetryPolicy) {
	ctx, cancel := context.WithTimeout(ctx, cfg.Timeout)
	defer cancel()
	return op(ctx)
}

func withRetry[T any](cfg *QueryConfig, op func() (_res *T, _err *PgLibError, _shouldRetry RetryPolicy)) (*T, *PgLibError) {
	var lastErr *PgLibError = nil
	maxAllowedFree := 5
	for i := cfg.Retries; i > 0; i -= 1 {
		res, err, shouldRetry := op()
		if err == nil {
			return res, nil
		}
		if shouldRetry == NoRetry {
			return nil, err
		}
		if shouldRetry == FreeRetry && maxAllowedFree > 0 {
			i += 1
			maxAllowedFree -= 1
		}
		lastErr = err
	}
	if lastErr == nil {
		return nil, makeServerError(fmt.Errorf("no operation was performed"))
	}
	return nil, makeConnectionError(lastErr)
}

func PerformOperation[T any](ctx context.Context, pool *pgxpool.Pool, cfg *QueryConfig, op func(context.Context, *pgxpool.Pool) (*T, error, RetryPolicy)) (*T, *PgLibError) {
	res, err := withRetry[T](cfg, func() (*T, *PgLibError, RetryPolicy) {
		return withTimeout[T](ctx, cfg, func(timeoutCtx context.Context) (*T, *PgLibError, RetryPolicy) {
			return withTransaction[T](timeoutCtx, pool, cfg, op)
		})
	})
	return res, err
}
