package tracectx

import (
	"context"
)

type (
	traceIDKey struct{} // store traceID in context
)

func NewTraceIDContext(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey{}, traceID)
}

// FromTraceIDContext get trace id from context.
func FromTraceIDContext(ctx context.Context) string {
	if v := ctx.Value(traceIDKey{}); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func WithTraceIDContext(ctx context.Context) context.Context {
	if v := ctx.Value(traceIDKey{}); v == nil {
		traceID := NewTraceID()
		return context.WithValue(ctx, traceIDKey{}, traceID)
	}
	return ctx
}
