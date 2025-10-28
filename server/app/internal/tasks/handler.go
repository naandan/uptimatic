package tasks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"uptimatic/internal/adapters/email"
	"uptimatic/internal/config"
	"uptimatic/internal/models"
	"uptimatic/internal/url"
	"uptimatic/internal/utils"

	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

type TaskHandler struct {
	cfg      *config.Config
	pgsql    *gorm.DB
	client   *asynq.Client
	mailTask *email.EmailTask
	urlRepo  url.UrlRepository
	logRepo  url.StatusLogRepository
}

func NewTaskHandler(cfg *config.Config, pgsql *gorm.DB, client *asynq.Client, mailTask *email.EmailTask, urlRepo url.UrlRepository, logRepo url.StatusLogRepository) *TaskHandler {
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

func (h *TaskHandler) enqueueEmail(to, subject string, emailType email.EmailType, data map[string]any) error {
	emailPayload, err := json.Marshal(email.EmailPayload{
		To:      to,
		Subject: subject,
		Type:    emailType,
		Data:    data,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal email payload: %w", err)
	}

	task := asynq.NewTask(TaskSendEmail, emailPayload)
	if _, err := h.client.Enqueue(task); err != nil {
		return fmt.Errorf("failed to enqueue email task: %w", err)
	}
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

	lastLog, err := h.logRepo.GetLastLogByURLID(ctx, h.pgsql, payload.ID)
	var lastStatus int64
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Debug(ctx, "No previous log found, treating as first check", nil)
			lastStatus = 0
		} else {
			utils.Error(ctx, "Failed to get last log", map[string]any{"error": err.Error()})
			return fmt.Errorf("failed to get last log: %w", err)
		}
	} else {
		lastStatus, err = strconv.ParseInt(lastLog.Status, 10, 64)
		if err != nil {
			utils.Error(ctx, "Failed to parse last status", map[string]any{"error": err.Error()})
			return fmt.Errorf("failed to parse last status: %w", err)
		}
	}

	client := &http.Client{Timeout: 30 * time.Second}
	ctxReq, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctxReq, http.MethodGet, payload.URL, nil)
	if err != nil {
		utils.Error(ctx, "Failed to create HTTP request", map[string]any{"error": err.Error()})
		return fmt.Errorf("failed to create request: %w", err)
	}

	start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		utils.Error(ctx, "HTTP request failed", map[string]any{"url": payload.URL, "error": err.Error()})
		return fmt.Errorf("failed to check URL %s: %w", payload.URL, err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			utils.Error(ctx, "Failed to close response body", map[string]any{"error": err.Error()})
		}
	}()

	duration := time.Since(start)

	log := models.StatusLog{
		URLID:        payload.ID,
		Status:       strconv.Itoa(resp.StatusCode),
		ResponseTime: int64(duration.Milliseconds()),
		CheckedAt:    time.Now().UTC(),
	}

	if err := h.logRepo.Create(ctx, h.pgsql, &log); err != nil {
		utils.Error(ctx, "Failed to create status log", map[string]any{"url_id": payload.ID, "error": err.Error()})
		return fmt.Errorf("failed to create status log: %w", err)
	}

	utils.Debug(ctx, "URL checked result", map[string]any{
		"url":           payload.URL,
		"status":        log.Status,
		"response_time": log.ResponseTime,
	})

	const downThreshold = 400
	if resp.StatusCode != int(lastStatus) {
		loc, _ := time.LoadLocation("Asia/Jakarta")

		data := map[string]any{
			"LogoURL":      fmt.Sprintf("%s://%s/icon.png", h.cfg.AppScheme, h.cfg.AppDomain),
			"Label":        payload.Label,
			"URL":          payload.URL,
			"Status":       log.Status,
			"ResponseTime": log.ResponseTime,
			"CheckedAt":    log.CheckedAt.In(loc).Format("2006-01-02 15:04:05"),
		}

		if resp.StatusCode >= downThreshold {
			utils.Warn(ctx, "URL is down, sending notification", map[string]any{
				"url":    payload.URL,
				"status": resp.StatusCode,
			})

			if err := h.enqueueEmail(payload.User.Email, "Uptime Alert - Website Down", email.EmailDown, data); err != nil {
				utils.Error(ctx, "Failed to enqueue down email", map[string]any{"error": err.Error()})
				return fmt.Errorf("failed to enqueue down email: %w", err)
			}
		} else {
			utils.Info(ctx, "URL is up, sending notification", map[string]any{
				"url":    payload.URL,
				"status": resp.StatusCode,
			})

			if err := h.enqueueEmail(payload.User.Email, "Uptime Alert - Website Up", email.EmailUp, data); err != nil {
				utils.Error(ctx, "Failed to enqueue up email", map[string]any{"error": err.Error()})
				return fmt.Errorf("failed to enqueue up email: %w", err)
			}
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
