package utils

import (
	"context"

	"github.com/google/uuid"
)

type ctxKey string

const TraceKey ctxKey = "trace_id"

func WithTraceID(ctx context.Context) context.Context {
	traceID := uuid.NewString()
	return context.WithValue(ctx, TraceKey, traceID)
}
