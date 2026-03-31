package include

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"team_dynamics/pg_lib/src/pg_lib"
)

type (
	ConnectionConfig     = pg_lib.ConnectionConfig
	InitializationConfig = pg_lib.InitializationConfig
	PoolConfig           = pg_lib.PoolConfig
	QueryConfig          = pg_lib.QueryConfig
)

var (
	MakePool = pg_lib.MakePool
)

func PerformOperation[T any](ctx context.Context, pool *pgxpool.Pool, cfg *QueryConfig, op func(*pgxpool.Pool) (*T, error)) (*T, error) {
	return pg_lib.PerformOperation(ctx, pool, cfg, op)
}
