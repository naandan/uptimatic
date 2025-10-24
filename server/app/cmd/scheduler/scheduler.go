package scheduler

import (
	"context"
	"time"
	"uptimatic/internal/config"
	"uptimatic/internal/db"
	"uptimatic/internal/tasks"
	"uptimatic/internal/utils"

	"github.com/getsentry/sentry-go"
	"github.com/hibiken/asynq"
)

func Start() {
	ctx := context.Background()
	ctx = utils.WithTraceID(ctx)

	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	utils.InitLogger(cfg.AppLogLevel)
	_ = utils.InitSentry(cfg.SentryDSN)
	defer sentry.Flush(2 * time.Second)

	scheduler := db.NewAsynqScheduler(&cfg)

	_, err = scheduler.Register(
		"* * * * *",
		asynq.NewTask(tasks.TaskValidateUptime, nil),
	)
	if err != nil {
		utils.Fatal(ctx, "Failed to register task", map[string]any{"error": err})
		return
	}

	utils.Debug(ctx, "Scheduler started", nil)
	if err := scheduler.Run(); err != nil {
		utils.Fatal(ctx, "Failed to run scheduler", map[string]any{"error": err})
	}
}
