package worker

import (
	"context"
	"time"
	"uptimatic/internal/adapters/email"
	"uptimatic/internal/config"
	"uptimatic/internal/db"
	"uptimatic/internal/tasks"
	"uptimatic/internal/url"
	"uptimatic/internal/utils"

	"github.com/getsentry/sentry-go"
	"github.com/hibiken/asynq"
)

func Start() {
	ctx := context.Background()
	ctx = utils.WithTraceID(ctx)

	cfg, err := config.LoadConfig()
	if err != nil {
		utils.Fatal(ctx, "Failed to load config", map[string]any{"error": err})
	}

	utils.InitLogger(cfg.AppLogLevel)
	_ = utils.InitSentry(cfg.SentryDSN)
	defer sentry.Flush(2 * time.Second)

	psql := db.NewPostgresClient(&cfg)
	client := db.NewAsynqClient(&cfg)

	urlRepo := url.NewUrlRepository()
	logRepo := url.NewLogRepository()

	mailTask, err := email.NewEmailTask(&cfg)
	if err != nil {
		utils.Fatal(ctx, "Failed to create email task", map[string]any{"error": err})
	}

	handler := tasks.NewTaskHandler(&cfg, psql, client, mailTask, urlRepo, logRepo)

	srv := db.NewAsynqServer(&cfg)
	mux := asynq.NewServeMux()

	mux.HandleFunc(tasks.TaskSendEmail, tasks.MiddlewareHandler(handler.SendEmailHandler))
	mux.HandleFunc(tasks.TaskValidateUptime, tasks.MiddlewareHandler(handler.ValidateUptimeHandler))
	mux.HandleFunc(tasks.TaskCheckUptime, tasks.MiddlewareHandler(handler.CheckUptimeHandler))

	utils.Debug(ctx, "Worker started", nil)
	if err := srv.Run(mux); err != nil {
		utils.Fatal(ctx, "Failed to start server", map[string]any{"error": err})
	}
}
