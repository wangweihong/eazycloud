package log

import (
	"context"
)

type logKeyCtx struct{}

// WithContext returns a copy of context in which the log value is set.
func WithContext(ctx context.Context) context.Context {
	return std.WithContext(ctx)
}

// save log handler into zap
func (l *zapLogger) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, logKeyCtx{}, l)
}

// FromContext returns the value of the log key on the ctx.
func FromContext(ctx context.Context) Logger {
	if ctx != nil {
		logger := ctx.Value(logKeyCtx{})
		if logger != nil {
			return logger.(Logger)
		}
	}

	return WithName("Unknown-Context")
}

type fieldKeyCtx struct{}

// WithContext returns a copy of context in which the log value is set.
func WithFields(ctx context.Context, fields map[string]interface{}) context.Context {
	return context.WithValue(ctx, fieldKeyCtx{}, fields)
}
