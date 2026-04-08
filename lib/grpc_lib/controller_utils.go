package grpc_lib

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log/slog"
	"logging"
	"os"
)

func HandlePanic(ctx context.Context, errOut *error) {
	logger := logging.GetLogger(ctx)
	if r := recover(); r != nil {
		msg := fmt.Sprintf("panic occured: %v", r)
		logger.Error(msg)
		*errOut = status.Error(codes.Internal, msg)
	}
}

func WithLogger(ctx context.Context, rpc string) context.Context {
	requestId := uuid.New().String()
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		requestIdOpt := md.Get("request-id")
		if len(requestIdOpt) > 0 {
			requestId = requestIdOpt[0]
		}
	}
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug})
	logger := slog.New(handler).With("requestId", requestId, "rpc", rpc)
	return logging.WithLogger(ctx, logger)
}
