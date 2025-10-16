package worker

import (
	"context"
	"uptimatic/internal/config"
	"uptimatic/internal/db"
	"uptimatic/internal/email"
	"uptimatic/internal/repositories"
	"uptimatic/internal/tasks"
	"uptimatic/internal/utils"

	"github.com/hibiken/asynq"
)

func Start() {
	ctx := context.Background()
	ctx = utils.WithTraceID(ctx)

	cfg, err := config.LoadConfig()
	if err != nil {
		utils.Fatal(ctx, "failed to load config", map[string]any{"error": err})
	}

	utils.InitLogger(cfg.AppLogLevel)

	psql := db.NewPostgresClient(&cfg)
	client := db.NewAsynqClient(&cfg)

	urlRepo := repositories.NewUrlRepository()
	logRepo := repositories.NewLogRepository()

	mailTask, err := email.NewEmailTask(&cfg)
	if err != nil {
		utils.Fatal(ctx, "failed to create email task", map[string]any{"error": err})
	}

	handler := tasks.NewTaskHandler(&cfg, psql, client, mailTask, urlRepo, logRepo)

	srv := db.NewAsynqServer(&cfg)
	mux := asynq.NewServeMux()

	mux.HandleFunc(tasks.TaskSendEmail, tasks.MiddlewareHandler(handler.SendEmailHandler))
	mux.HandleFunc(tasks.TaskValidateUptime, tasks.MiddlewareHandler(handler.ValidateUptimeHandler))
	mux.HandleFunc(tasks.TaskCheckUptime, tasks.MiddlewareHandler(handler.CheckUptimeHandler))

	utils.Debug(ctx, "Worker started", nil)
	if err := srv.Run(mux); err != nil {
		utils.Fatal(ctx, "failed to start server", map[string]any{"error": err})
	}
}
