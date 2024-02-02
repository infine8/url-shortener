package grpc

import (
	"context"
	"log/slog"

	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

func interceptorLogger(l *slog.Logger) grpclog.Logger {
    return grpclog.LoggerFunc(func(ctx context.Context, lvl grpclog.Level, msg string, fields ...any) {
        l.Log(ctx, slog.Level(lvl), msg, fields...)
    })
}