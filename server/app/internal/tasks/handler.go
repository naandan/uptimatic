package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"uptimatic/internal/config"
	"uptimatic/internal/email"
	"uptimatic/internal/models"
	"uptimatic/internal/repositories"
	"uptimatic/internal/utils"

	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

type TaskHandler struct {
	cfg      *config.Config
	pgsql    *gorm.DB
	client   *asynq.Client
	mailTask *email.EmailTask
	urlRepo  repositories.UrlRepository
	logRepo  repositories.StatusLogRepository
}

func NewTaskHandler(cfg *config.Config, pgsql *gorm.DB, client *asynq.Client, mailTask *email.EmailTask, urlRepo repositories.UrlRepository, logRepo repositories.StatusLogRepository) *TaskHandler {
	return &TaskHandler{cfg, pgsql, client, mailTask, urlRepo, logRepo}
}

func (h *TaskHandler) SendEmailHandler(ctx context.Context, t *asynq.Task) error {
	var payload email.EmailPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}
	utils.Info(nil, "Sending email", map[string]any{"to": payload.To, "subject": payload.Subject})
	return h.mailTask.SendEmail(ctx, payload.To, payload.Subject, payload.Type, payload.Data)
}

func (h *TaskHandler) CheckUptimeHandler(ctx context.Context, t *asynq.Task) error {
	var payload models.URL
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	client := http.Client{Timeout: 30 * time.Second}

	start := time.Now()
	resp, err := client.Get(payload.URL)
	duration := time.Since(start)
	if err != nil {
		return fmt.Errorf("failed to check URL %s: %w", payload.URL, err)
	}
	defer resp.Body.Close()

	log := models.StatusLog{
		URLID:        payload.ID,
		Status:       strconv.Itoa(resp.StatusCode),
		ResponseTime: int64(duration.Milliseconds()),
		CheckedAt:    time.Now().UTC(),
	}

	if err := h.logRepo.Create(h.pgsql, &log); err != nil {
		return fmt.Errorf("failed to create status log: %w", err)
	}

	payload.LastChecked = &log.CheckedAt
	if err := h.urlRepo.Update(h.pgsql, &payload); err != nil {
		return fmt.Errorf("failed to update URL: %w", err)
	}

	utils.Debug(nil, "URL checked", map[string]any{
		"url":           payload.URL,
		"status":        log.Status,
		"response_time": log.ResponseTime,
	})
	return nil
}

func (h *TaskHandler) ValidateUptimeHandler(ctx context.Context, t *asynq.Task) error {
	urls, err := h.urlRepo.GetActiveURLs(h.pgsql)
	if err != nil {
		return fmt.Errorf("failed to get active URLs: %w", err)
	}

	now := time.Now().UTC()

	for _, url := range urls {
		if url.LastChecked == nil {
			payload, err := json.Marshal(url)
			if err != nil {
				return fmt.Errorf("failed to marshal URL payload: %w", err)
			}

			task := asynq.NewTask(TaskCheckUptime, payload)
			if _, err := h.client.Enqueue(task); err != nil {
				utils.Error(nil, "failed to enqueue first check task", map[string]any{"url": url.URL, "error": err})
				continue
			}

			utils.Debug(nil, "First uptime check scheduled", map[string]any{"url": url.URL})
			continue
		}

		diff := now.Sub(url.LastChecked.UTC()) + 5

		if diff >= time.Duration(url.Interval)*time.Second {
			payload, err := json.Marshal(url)
			if err != nil {
				return fmt.Errorf("failed to marshal URL payload: %w", err)
			}

			task := asynq.NewTask(TaskCheckUptime, payload)
			if _, err := h.client.Enqueue(task); err != nil {
				utils.Error(nil, "failed to enqueue uptime check", map[string]any{"url": url.URL, "error": err})
				continue
			}

			utils.Debug(nil, "URL needs to be checked", map[string]any{"url": url.URL})
		}
	}

	return nil
}
