package repositories

import (
	"context"
	"uptimatic/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UrlRepository interface {
	Create(ctx context.Context, tx *gorm.DB, url *models.URL) error
	Update(ctx context.Context, tx *gorm.DB, url *models.URL) error
	Delete(ctx context.Context, tx *gorm.DB, url *models.URL) error
	FindByPublicID(ctx context.Context, tx *gorm.DB, publicID uuid.UUID) (*models.URL, error)
	ListByUserID(ctx context.Context, tx *gorm.DB, userID uint, page, perPage int, active *bool, searchLabel string, sortBy string) ([]models.URL, int, error)
	GetActiveURLs(ctx context.Context, tx *gorm.DB) ([]models.URL, error)
}

type urlRepository struct{}

func NewUrlRepository() UrlRepository {
	return &urlRepository{}
}

func (r *urlRepository) Create(ctx context.Context, tx *gorm.DB, url *models.URL) error {
	return tx.WithContext(ctx).Create(url).Error
}

func (r *urlRepository) Update(ctx context.Context, tx *gorm.DB, url *models.URL) error {
	return tx.WithContext(ctx).Save(url).Error
}

func (r *urlRepository) Delete(ctx context.Context, tx *gorm.DB, url *models.URL) error {
	return tx.WithContext(ctx).Delete(url).Error
}

func (r *urlRepository) FindByPublicID(ctx context.Context, tx *gorm.DB, publicID uuid.UUID) (*models.URL, error) {
	var url models.URL
	err := tx.WithContext(ctx).First(&url, "public_id = ?", publicID).Error
	if err != nil {
		return nil, err
	}
	return &url, nil
}

func (r *urlRepository) ListByUserID(
	ctx context.Context,
	tx *gorm.DB,
	userID uint,
	page, perPage int,
	active *bool,
	searchLabel string,
	sortBy string,
) ([]models.URL, int, error) {

	var urls []models.URL
	var count int64

	query := tx.WithContext(ctx).Model(&models.URL{}).Where("user_id = ?", userID)

	if active != nil {
		query = query.Where("active = ?", *active)
	}

	if searchLabel != "" {
		query = query.Where("label ILIKE ?", "%"+searchLabel+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if sortBy == "created_at" {
		query = query.Order("created_at " + "DESC")
	} else {
		query = query.Order("label" + " " + "ASC")
	}

	if err := query.Offset((page - 1) * perPage).Limit(perPage).Find(&urls).Error; err != nil {
		return nil, 0, err
	}

	return urls, int(count), nil
}

func (r *urlRepository) GetActiveURLs(ctx context.Context, tx *gorm.DB) ([]models.URL, error) {
	var urls []models.URL
	err := tx.WithContext(ctx).Preload("User").Where("active = ?", true).Find(&urls).Error
	if err != nil {
		return nil, err
	}
	return urls, nil
}
