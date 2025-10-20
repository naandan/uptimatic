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
		utils.Error(ctx, "Failed to unmarshal email payload", map[string]any{"error": err.Error()})
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	utils.Info(ctx, "Sending email", map[string]any{
		"to":      payload.To,
		"subject": payload.Subject,
		"type":    payload.Type,
	})

	if err := h.mailTask.SendEmail(ctx, payload.To, payload.Subject, payload.Type, payload.Data); err != nil {
		utils.Error(ctx, "Failed to send email", map[string]any{
			"to":    payload.To,
			"error": err.Error(),
		})
		return err
	}

	utils.Info(ctx, "Email sent successfully", map[string]any{
		"to":      payload.To,
		"subject": payload.Subject,
	})
	return nil
}

func (h *TaskHandler) CheckUptimeHandler(ctx context.Context, t *asynq.Task) error {
	var payload models.URL
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		utils.Error(ctx, "Failed to unmarshal uptime payload", map[string]any{"error": err.Error()})
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	utils.Info(ctx, "Checking URL uptime", map[string]any{
		"url_id": payload.ID,
		"url":    payload.URL,
		"label":  payload.Label,
	})

	client := http.Client{Timeout: 30 * time.Second}
	start := time.Now()
	resp, err := client.Get(payload.URL)
	duration := time.Since(start)

	if err != nil {
		utils.Error(ctx, "HTTP request failed", map[string]any{
			"url":   payload.URL,
			"error": err.Error(),
		})
		return fmt.Errorf("failed to check URL %s: %w", payload.URL, err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			utils.Error(ctx, "Failed to close response body", map[string]any{"error": err.Error()})
		}
	}()

	log := models.StatusLog{
		URLID:        payload.ID,
		Status:       strconv.Itoa(resp.StatusCode),
		ResponseTime: int64(duration.Milliseconds()),
		CheckedAt:    time.Now().UTC(),
	}

	if err := h.logRepo.Create(ctx, h.pgsql, &log); err != nil {
		utils.Error(ctx, "Failed to create status log", map[string]any{
			"url_id": payload.ID,
			"error":  err.Error(),
		})
		return fmt.Errorf("failed to create status log: %w", err)
	}

	utils.Debug(ctx, "URL checked result", map[string]any{
		"url":           payload.URL,
		"status":        log.Status,
		"response_time": log.ResponseTime,
	})

	if resp.StatusCode >= 400 {
		utils.Warn(ctx, "URL is down, sending notification", map[string]any{
			"url":    payload.URL,
			"status": resp.StatusCode,
		})

		loc, _ := time.LoadLocation("Asia/Jakarta")
		emailPayload, err := json.Marshal(email.EmailPayload{
			To:      payload.User.Email,
			Subject: "Uptime Alert - Website Down",
			Type:    email.EmailDown,
			Data: map[string]any{
				"LogoURL":      fmt.Sprintf("%s://%s/icon.png", h.cfg.AppScheme, h.cfg.AppDomain),
				"Label":        payload.Label,
				"URL":          payload.URL,
				"Status":       log.Status,
				"ResponseTime": log.ResponseTime,
				"CheckedAt":    log.CheckedAt.In(loc).Format("2006-01-02 15:04:05"),
			},
		})

		if err != nil {
			utils.Error(ctx, "Failed to encode email payload", map[string]any{"error": err.Error()})
			return fmt.Errorf("failed to json encode payload: %w", err)
		}

		task := asynq.NewTask(TaskSendEmail, emailPayload)
		if _, err := h.client.Enqueue(task); err != nil {
			utils.Error(ctx, "Failed to enqueue email task", map[string]any{"url": payload.URL, "error": err.Error()})
			return fmt.Errorf("failed to enqueue email task: %w", err)
		}
	}

	payload.LastChecked = &log.CheckedAt
	if err := h.urlRepo.Update(ctx, h.pgsql, &payload); err != nil {
		utils.Error(ctx, "Failed to update URL last checked", map[string]any{
			"url_id": payload.ID,
			"error":  err.Error(),
		})
		return fmt.Errorf("failed to update URL: %w", err)
	}

	utils.Info(ctx, "Uptime check completed successfully", map[string]any{
		"url_id": payload.ID,
		"url":    payload.URL,
	})
	return nil
}

func (h *TaskHandler) ValidateUptimeHandler(ctx context.Context, t *asynq.Task) error {
	utils.Info(ctx, "Running uptime validation task", nil)

	urls, err := h.urlRepo.GetActiveURLs(ctx, h.pgsql)
	if err != nil {
		utils.Error(ctx, "Failed to get active URLs", map[string]any{"error": err.Error()})
		return fmt.Errorf("failed to get active URLs: %w", err)
	}

	now := time.Now().UTC()
	utils.Debug(ctx, "Active URLs fetched", map[string]any{"count": len(urls)})

	for _, url := range urls {
		if url.LastChecked == nil {
			utils.Info(ctx, "Scheduling first uptime check", map[string]any{"url": url.URL})

			payload, err := json.Marshal(url)
			if err != nil {
				utils.Error(ctx, "Failed to marshal first-check URL payload", map[string]any{"url": url.URL, "error": err.Error()})
				continue
			}

			task := asynq.NewTask(TaskCheckUptime, payload)
			if _, err := h.client.Enqueue(task); err != nil {
				utils.Error(ctx, "Failed to enqueue first uptime check task", map[string]any{"url": url.URL, "error": err.Error()})
				continue
			}
			continue
		}

		diff := now.Sub(url.LastChecked.UTC())
		if diff >= time.Duration(url.Interval)*time.Second {
			utils.Debug(ctx, "URL due for next check", map[string]any{
				"url":       url.URL,
				"interval":  url.Interval,
				"lastCheck": url.LastChecked,
			})

			payload, err := json.Marshal(url)
			if err != nil {
				utils.Error(ctx, "Failed to marshal URL payload for recheck", map[string]any{"url": url.URL, "error": err.Error()})
				continue
			}

			task := asynq.NewTask(TaskCheckUptime, payload)
			if _, err := h.client.Enqueue(task); err != nil {
				utils.Error(ctx, "Failed to enqueue uptime check", map[string]any{"url": url.URL, "error": err.Error()})
				continue
			}
		}
	}

	utils.Info(ctx, "Uptime validation task completed", map[string]any{"total_urls": len(urls)})
	return nil
}
