package pg_lib

import (
	"github.com/jackc/pgx/v5"
	"time"
)

type ConnectionConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSLMode  string
}

type InitializationConfig struct {
	ConnectionTimeout time.Duration
	Retries           int32
}

type PoolConfig struct {
	MaxConns          int32
	MinConns          int32
	MaxConnLifetime   time.Duration
	MaxConnIdleTime   time.Duration
	HealthCheckPeriod time.Duration
}

type QueryConfig struct {
	Retries        int32
	Timeout        time.Duration
	IsolationLevel pgx.TxIsoLevel
	AccessMode     pgx.TxAccessMode
}
