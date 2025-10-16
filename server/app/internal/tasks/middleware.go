package tasks

import (
	"context"
	"uptimatic/internal/utils"

	"github.com/hibiken/asynq"
)

// MiddlewareHandler membungkus handler agar otomatis inject trace_id.
func MiddlewareHandler(handler func(context.Context, *asynq.Task) error) asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		ctx = utils.WithTraceID(ctx)
		utils.Debug(ctx, "Task started", map[string]any{"type": t.Type()})
		err := handler(ctx, t)
		if err != nil {
			utils.Error(ctx, "Task failed", map[string]any{"error": err.Error(), "type": t.Type()})
		} else {
			utils.Debug(ctx, "Task completed", map[string]any{"type": t.Type()})
		}
		return err
	}
}
