package pg_lib

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

func makeConnectionString(cfg *ConnectionConfig) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database, cfg.SSLMode,
	)
}

func MakePool(ctx context.Context, connCfg *ConnectionConfig, poolCfg *PoolConfig, initCfg *InitializationConfig) (*pgxpool.Pool, error) {
	dsn := makeConnectionString(connCfg)
	var res *pgxpool.Pool
	for i := initCfg.Retries; i > 0; i -= 1 {
		config, err := pgxpool.ParseConfig(dsn)
		if err != nil {
			continue
		}
		config.MaxConns = poolCfg.MaxConns
		config.MaxConns = poolCfg.MinConns
		config.MaxConnLifetime = poolCfg.MaxConnLifetime
		config.MaxConnIdleTime = poolCfg.MaxConnIdleTime
		config.HealthCheckPeriod = poolCfg.HealthCheckPeriod

		config.ConnConfig.ConnectTimeout = initCfg.ConnectionTimeout
		config.ConnConfig.RuntimeParams["timezone"] = "UTC"

		pool, err := pgxpool.NewWithConfig(ctx, config)
		if err != nil {
			continue
		}
		res = pool
	}
	if res == nil {
		return nil, fmt.Errorf("cannot initialize pg connection")
	}
	return res, nil
}
