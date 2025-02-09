package logging

import (
	"context"

	"go.uber.org/zap"
)

type loggerKey struct{}

func WithContext(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

func FromContext(ctx context.Context) *zap.Logger {
	if ctx == nil {
		panic("nil logger")
	}

	if logger, _ := ctx.Value(loggerKey{}).(*zap.Logger); logger != nil {
		return logger
	}

	panic("no logger in context")
}
