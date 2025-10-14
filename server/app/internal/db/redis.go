package db

import (
	"context"
	"fmt"
	"uptimatic/internal/config"
	"uptimatic/internal/utils"

	"github.com/go-redis/redis/v8"
	"github.com/hibiken/asynq"
)

func NewRedisClient(cfg *config.Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisHost + ":" + fmt.Sprint(cfg.RedisPort),
		Password: cfg.RedisPass,
		DB:       0,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		utils.Fatal(nil, "failed to connect redis", map[string]any{"error": err})
	}

	return rdb
}

func NewAsynqClient(cfg *config.Config) *asynq.Client {
	return asynq.NewClient(asynq.RedisClientOpt{Addr: cfg.RedisHost + ":" + fmt.Sprint(cfg.RedisPort)})
}

func NewAsynqServer(cfg *config.Config) *asynq.Server {
	return asynq.NewServer(
		asynq.RedisClientOpt{Addr: cfg.RedisHost + ":" + fmt.Sprint(cfg.RedisPort)},
		asynq.Config{
			Concurrency: 10,
		},
	)
}

func NewAsynqScheduler(cfg *config.Config) *asynq.Scheduler {
	return asynq.NewScheduler(
		asynq.RedisClientOpt{Addr: cfg.RedisHost + ":" + fmt.Sprint(cfg.RedisPort)},
		&asynq.SchedulerOpts{},
	)
}

func RedisClientOpt(cfg *config.Config) asynq.RedisClientOpt {
	return asynq.RedisClientOpt{Addr: cfg.RedisHost + ":" + fmt.Sprint(cfg.RedisPort)}
}
