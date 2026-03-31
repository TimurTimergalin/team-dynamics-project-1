package pg_lib

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func withTransaction[T any](ctx context.Context, pool *pgxpool.Pool, cfg *QueryConfig, op func(*pgxpool.Pool) (*T, error)) (*T, error) {
	tx, err := pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   cfg.IsolationLevel,
		AccessMode: cfg.AccessMode,
	})
	if err != nil {
		return nil, err
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		_ = tx.Rollback(ctx)
	}(tx, ctx)

	res, err := op(pool)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func withTimeout[T any](ctx context.Context, cfg *QueryConfig, op func(context.Context) (*T, error)) (*T, error) {
	ctx, cancel := context.WithTimeout(ctx, cfg.Timeout)
	defer cancel()
	return op(ctx)
}

func withRetry[T any](cfg *QueryConfig, op func() (*T, error)) (*T, error) {
	var lastErr error = nil
	for i := cfg.Retries; i > 0; i -= 1 {
		res, err := op()
		if err == nil {
			return res, nil
		}
		lastErr = err
	}
	if lastErr == nil {
		return nil, fmt.Errorf("no tries were performed")
	}
	return nil, fmt.Errorf("max retiries exceeded. Last error: %w", lastErr)
}

func PerformOperation[T any](ctx context.Context, pool *pgxpool.Pool, cfg *QueryConfig, op func(*pgxpool.Pool) (*T, error)) (*T, error) {
	res, err := withRetry[T](cfg, func() (*T, error) {
		return withTimeout[T](ctx, cfg, func(timeoutCtx context.Context) (*T, error) {
			return withTransaction[T](timeoutCtx, pool, cfg, op)
		})
	})
	return res, err
}
