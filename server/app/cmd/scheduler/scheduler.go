package scheduler

import (
	"uptimatic/internal/config"
	"uptimatic/internal/db"
	"uptimatic/internal/tasks"
	"uptimatic/internal/utils"

	"github.com/hibiken/asynq"
)

func Start() {
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
		utils.Fatal(nil, "failed to register task", map[string]any{"error": err})
		return
	}

	utils.Debug(nil, "Scheduler started", nil)
	scheduler.Run()
}
