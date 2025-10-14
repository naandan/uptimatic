package worker

import (
	"uptimatic/internal/config"
	"uptimatic/internal/db"
	"uptimatic/internal/email"
	"uptimatic/internal/repositories"
	"uptimatic/internal/tasks"
	"uptimatic/internal/utils"

	"github.com/hibiken/asynq"
)

func Start() {
	cfg, err := config.LoadConfig()
	if err != nil {
		utils.Fatal(nil, "failed to load config", map[string]any{"error": err})
	}

	utils.InitLogger(cfg.AppLogLevel)

	psql := db.NewPostgresClient(&cfg)
	client := db.NewAsynqClient(&cfg)

	urlRepo := repositories.NewUrlRepository()
	logRepo := repositories.NewLogRepository()

	mailTask, err := email.NewEmailTask(&cfg)
	if err != nil {
		utils.Fatal(nil, "failed to create email task", map[string]any{"error": err})
	}

	handler := tasks.NewTaskHandler(&cfg, psql, client, mailTask, urlRepo, logRepo)

	srv := db.NewAsynqServer(&cfg)
	mux := asynq.NewServeMux()

	mux.HandleFunc(tasks.TaskSendEmail, handler.SendEmailHandler)
	mux.HandleFunc(tasks.TaskValidateUptime, handler.ValidateUptimeHandler)
	mux.HandleFunc(tasks.TaskCheckUptime, handler.CheckUptimeHandler)

	utils.Debug(nil, "Worker started", nil)
	if err := srv.Run(mux); err != nil {
		utils.Fatal(nil, "failed to start server", map[string]any{"error": err})
	}
}
