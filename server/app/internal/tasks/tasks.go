package tasks

import (
	"encoding/json"
	"uptimatic/internal/adapters/email"

	"github.com/hibiken/asynq"
)

const (
	TaskSendEmail      = "send_email"
	TaskValidateUptime = "validate_uptime"
	TaskCheckUptime    = "check_uptime"
)

func EnqueueEmail(client *asynq.Client, to, subject string, mailType email.EmailType, data map[string]any) error {
	payload, err := json.Marshal(email.EmailPayload{
		To:      to,
		Subject: subject,
		Type:    mailType,
		Data:    data,
	})
	if err != nil {
		return err
	}

	task := asynq.NewTask(TaskSendEmail, payload)
	_, err = client.Enqueue(task,
		asynq.MaxRetry(3),
		asynq.ProcessIn(0),
	)
	return err
}
