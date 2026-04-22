package include

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"team_dynamics/pg_lib/src/pg_lib"
)

type (
	ConnectionConfig     = pg_lib.ConnectionConfig
	InitializationConfig = pg_lib.InitializationConfig
	PoolConfig           = pg_lib.PoolConfig
	QueryConfig          = pg_lib.QueryConfig
	ResponseStatus       = pg_lib.ResponseStatus
	PgLibError           = pg_lib.PgLibError
	PgLibErrorType       = pg_lib.PgLibErrorType
)

var (
	MakePool             = pg_lib.MakePool
	IsNoRows             = pg_lib.IsNoRows
	IsConstraintViolated = pg_lib.IsConstraintViolated
	IsSerializationError = pg_lib.IsSerializationError
	IsServerSideTimeout  = pg_lib.IsServerSideTimeout
)

const (
	NoRetry         = pg_lib.NoRetry
	NormalRetry     = pg_lib.NormalRetry
	FreeRetry       = pg_lib.FreeRetry
	ForceRollback   = pg_lib.ForceRollback
	LogicError      = pg_lib.LogicError
	ServerError     = pg_lib.ServerError
	ConnectionError = pg_lib.ConnectionError
)

func PerformOperation[T any](ctx context.Context, pool *pgxpool.Pool, cfg *QueryConfig, op func(context.Context, pgx.Tx) (T, error, ResponseStatus)) (T, *PgLibError) {
	return pg_lib.PerformOperation(ctx, pool, cfg, op)
}
