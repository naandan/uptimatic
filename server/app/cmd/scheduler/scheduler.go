package scheduler

import (
	"context"
	"uptimatic/internal/config"
	"uptimatic/internal/db"
	"uptimatic/internal/tasks"
	"uptimatic/internal/utils"

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

	scheduler := db.NewAsynqScheduler(&cfg)

	_, err = scheduler.Register(
		"* * * * *",
		asynq.NewTask(tasks.TaskValidateUptime, nil),
	)
	if err != nil {
		utils.Fatal(ctx, "failed to register task", map[string]any{"error": err})
		return
	}

	utils.Debug(ctx, "Scheduler started", nil)
	if err := scheduler.Run(); err != nil {
		utils.Fatal(ctx, "failed to run scheduler", map[string]any{"error": err})
	}
}
