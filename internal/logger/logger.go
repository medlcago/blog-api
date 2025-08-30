package logger

import (
	"context"
	"log/slog"
	"os"
)

type contextKey string

const (
	RequestIDKey contextKey = "request_id"
)

type Env string

const (
	EnvLocal Env = "local"
	EnvDev   Env = "dev"
	EnvProd  Env = "prod"
)

func New(env Env) *slog.Logger {
	var handler slog.Handler

	switch env {
	case EnvLocal:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		})
	case EnvDev:
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		})
	case EnvProd:
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelInfo,
		})
	default:
		// fallback
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		})
	}

	return slog.New(handler)
}

func WithRequestID(log *slog.Logger, reqID string) *slog.Logger {
	return log.With(
		slog.String(string(RequestIDKey), reqID),
	)
}

func FromCtx(ctx context.Context, log *slog.Logger) *slog.Logger {
	if v := ctx.Value(RequestIDKey); v != nil {
		if reqID, ok := v.(string); ok {
			return log.With(slog.String(string(RequestIDKey), reqID))
		}
	}
	return log
}

func Err(err error) slog.Attr {
	return slog.Any("error", err)
}
