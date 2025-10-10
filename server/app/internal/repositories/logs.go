package repositories

import (
	"uptimatic/internal/models"

	"gorm.io/gorm"
)

type StatusLogRepository interface {
	Create(tx *gorm.DB, log *models.StatusLog) error
	ListByURLID(tx *gorm.DB, urlID uint) ([]models.StatusLog, error)
}

type statusLogRepository struct{}

func NewLogRepository() StatusLogRepository {
	return &statusLogRepository{}
}

func (r *statusLogRepository) Create(tx *gorm.DB, log *models.StatusLog) error {
	return tx.Create(log).Error
}

func (r *statusLogRepository) ListByURLID(tx *gorm.DB, urlID uint) ([]models.StatusLog, error) {
	var logs []models.StatusLog
	err := tx.Where("url_id = ?", urlID).Find(&logs).Error
	if err != nil {
		return nil, err
	}
	return logs, nil
}
