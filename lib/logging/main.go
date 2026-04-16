package logging

import (
	"context"
	"log/slog"
)

type ctxKey int32

const (
	loggerKey ctxKey = iota
)

func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func GetLogger(ctx context.Context) *slog.Logger {
	return ctx.Value(loggerKey).(*slog.Logger)
}
