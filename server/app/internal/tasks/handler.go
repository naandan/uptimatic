package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"uptimatic/internal/config"
	"uptimatic/internal/email"
	"uptimatic/internal/utils"

	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

type TaskHandler struct {
	cfg      *config.Config
	pgsql    *gorm.DB
	mailTask *email.EmailTask
}

func NewTaskHandler(cfg *config.Config, pgsql *gorm.DB, mailTask *email.EmailTask) *TaskHandler {
	return &TaskHandler{cfg, pgsql, mailTask}
}

func (h *TaskHandler) SendEmailHandler(ctx context.Context, t *asynq.Task) error {
	var payload email.EmailPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}
	utils.Info(nil, "Sending email", map[string]any{"to": payload.To, "subject": payload.Subject})
	return h.mailTask.SendEmail(ctx, payload.To, payload.Subject, payload.Type, payload.Data)
}
