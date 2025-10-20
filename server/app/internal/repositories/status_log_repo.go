package repositories

import (
	"context"
	"time"
	"uptimatic/internal/models"

	"gorm.io/gorm"
)

type StatusLogRepository interface {
	Create(ctx context.Context, tx *gorm.DB, log *models.StatusLog) error
	GetByID(ctx context.Context, tx *gorm.DB, id uint) (*models.StatusLog, error)
	ListByURLID(ctx context.Context, tx *gorm.DB, urlID uint) ([]models.StatusLog, error)
	GetUptimeStats(ctx context.Context, tx *gorm.DB, urlID uint, truncUnit string, start, end time.Time) ([]models.UptimeStat, error)
}

type statusLogRepository struct{}

func NewLogRepository() StatusLogRepository {
	return &statusLogRepository{}
}

func (r *statusLogRepository) Create(ctx context.Context, tx *gorm.DB, log *models.StatusLog) error {
	return tx.WithContext(ctx).Create(log).Error
}

func (r *statusLogRepository) GetByID(ctx context.Context, tx *gorm.DB, id uint) (*models.StatusLog, error) {
	var log models.StatusLog
	err := tx.WithContext(ctx).First(&log, id).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

func (r *statusLogRepository) ListByURLID(ctx context.Context, tx *gorm.DB, urlID uint) ([]models.StatusLog, error) {
	var logs []models.StatusLog
	err := tx.WithContext(ctx).Where("url_id = ?", urlID).Find(&logs).Error
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (r *statusLogRepository) GetUptimeStats(ctx context.Context, tx *gorm.DB, urlID uint, truncUnit string, start, end time.Time) ([]models.UptimeStat, error) {
	var results []models.UptimeStat

	query := `
		SELECT
			date_trunc(?, checked_at) AS bucket_start,
			COUNT(*) AS total_checks,
			COUNT(*) FILTER (WHERE status BETWEEN 200 AND 299) AS up_checks,
			ROUND(
				COUNT(*) FILTER (WHERE status BETWEEN 200 AND 299) * 100.0 / COUNT(*),
				2
			) AS uptime_percent
		FROM status_logs
		WHERE url_id = ?
		  AND checked_at BETWEEN ? AND ?
		GROUP BY bucket_start
		ORDER BY bucket_start ASC;
	`

	if err := tx.WithContext(ctx).Raw(query, truncUnit, urlID, start, end).Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}
