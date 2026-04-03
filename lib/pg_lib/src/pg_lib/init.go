package pg_lib

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

func makeConnectionString(cfg *ConnectionConfig) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database, cfg.SSLMode,
	)
}

func ping(ctx context.Context, pool *pgxpool.Pool, timeout time.Duration) error {
	pingCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return fmt.Errorf("error while pinging: %w", err)
	}
	return nil
}

func MakePool(ctx context.Context, connCfg *ConnectionConfig, poolCfg *PoolConfig, initCfg *InitializationConfig) (*pgxpool.Pool, error) {
	dsn := makeConnectionString(connCfg)
	var res *pgxpool.Pool
	var lastErr error = nil
	for i := initCfg.Retries; i > 0; i -= 1 {
		config, err := pgxpool.ParseConfig(dsn)
		if err != nil {
			lastErr = err
			continue
		}
		config.MaxConns = poolCfg.MaxConns
		config.MinConns = poolCfg.MinConns
		config.MaxConnLifetime = poolCfg.MaxConnLifetime
		config.MaxConnIdleTime = poolCfg.MaxConnIdleTime
		config.HealthCheckPeriod = poolCfg.HealthCheckPeriod

		config.ConnConfig.ConnectTimeout = initCfg.ConnectionTimeout
		config.ConnConfig.RuntimeParams["timezone"] = "UTC"

		pool, err := pgxpool.NewWithConfig(ctx, config)
		if err != nil {
			lastErr = err
			continue
		}

		if err := ping(ctx, pool, initCfg.ConnectionTimeout); err != nil {
			lastErr = err
			continue
		}

		res = pool
		break
	}
	if res == nil {
		return nil, makeConnectionError(lastErr)
	}

	return res, nil
}
