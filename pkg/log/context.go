package log

import (
	"context"
)

type (
	logKeyCtx   struct{} // store logger in context
	fieldKeyCtx struct{} // store fields in context
)

// WithContext returns a copy of context in which the log value is set.
func WithContext(ctx context.Context) context.Context {
	return std.WithContext(ctx)
}

// save log handler into zap.
func (l *zapLogger) WithContext(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

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

// WithContext returns a copy of context in which the log value is set.
func WithFields(ctx context.Context, fields map[string]interface{}) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	return context.WithValue(ctx, fieldKeyCtx{}, fields)
}

// WithContext returns a copy of context in which the log value is set.
func WithFieldPair(ctx context.Context, key string, value interface{}) context.Context {
	if key == "" {
		return ctx
	}

	if ctx == nil {
		ctx = context.Background()
	}

	fieldMap := make(map[string]interface{})
	if fields := ctx.Value(fieldKeyCtx{}); fields != nil {
		var ok bool
		if fieldMap, ok = fields.(map[string]interface{}); ok {
			if fieldMap == nil {
				fieldMap = make(map[string]interface{})
			}
		}
	}
	fieldMap[key] = value

	return context.WithValue(ctx, fieldKeyCtx{}, fieldMap)
}
