package observability

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type traceIDKey struct{}

// TraceIDFromContext extracts the trace ID from context.
func TraceIDFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(traceIDKey{}).(string); ok {
		return v
	}
	return ""
}

// LoggingUnaryInterceptor provides structured logging with traceId for every gRPC call.
func LoggingUnaryInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		traceID := extractOrGenerateTraceID(ctx)
		ctx = context.WithValue(ctx, traceIDKey{}, traceID)

		start := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(start)

		code := status.Code(err)
		logger.LogAttrs(ctx, logLevel(err), "grpc request",
			slog.String("traceId", traceID),
			slog.String("method", info.FullMethod),
			slog.String("code", code.String()),
			slog.Duration("duration", duration),
		)

		return resp, err
	}
}

func extractOrGenerateTraceID(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if vals := md.Get("x-trace-id"); len(vals) > 0 && vals[0] != "" {
			return vals[0]
		}
	}
	return uuid.NewString()
}

func logLevel(err error) slog.Level {
	if err != nil {
		return slog.LevelError
	}
	return slog.LevelInfo
}
