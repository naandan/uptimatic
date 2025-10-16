package db

import (
	"context"
	"uptimatic/internal/config"
	"uptimatic/internal/utils"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresClient(cfg *config.Config) *gorm.DB {
	db, err := gorm.Open(postgres.Open(cfg.DBDSN()), &gorm.Config{})
	if err != nil {
		utils.Fatal(context.Background(), "failed to connect postgres", map[string]any{"error": err})
	}
	return db
}
