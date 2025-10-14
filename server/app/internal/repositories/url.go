package repositories

import (
	"uptimatic/internal/models"

	"gorm.io/gorm"
)

type UrlRepository interface {
	Create(tx *gorm.DB, url *models.URL) error
	Update(tx *gorm.DB, url *models.URL) error
	Delete(tx *gorm.DB, url *models.URL) error
	FindByID(tx *gorm.DB, id uint) (*models.URL, error)
	ListByUserID(tx *gorm.DB, userID uint, page, perPage int, active *bool, searchLabel string, sortBy string) ([]models.URL, int, error)
	GetActiveURLs(tx *gorm.DB) ([]models.URL, error)
}

type urlRepository struct{}

func NewUrlRepository() UrlRepository {
	return &urlRepository{}
}

func (r *urlRepository) Create(tx *gorm.DB, url *models.URL) error {
	return tx.Create(url).Error
}

func (r *urlRepository) Update(tx *gorm.DB, url *models.URL) error {
	return tx.Save(url).Error
}

func (r *urlRepository) Delete(tx *gorm.DB, url *models.URL) error {
	return tx.Delete(url).Error
}

func (r *urlRepository) FindByID(tx *gorm.DB, id uint) (*models.URL, error) {
	var url models.URL
	err := tx.First(&url, id).Error
	if err != nil {
		return nil, err
	}
	return &url, nil
}

func (r *urlRepository) ListByUserID(
	tx *gorm.DB,
	userID uint,
	page, perPage int,
	active *bool,
	searchLabel string,
	sortBy string,
) ([]models.URL, int, error) {

	var urls []models.URL
	var count int64

	query := tx.Model(&models.URL{}).Where("user_id = ?", userID)

	if active != nil {
		query = query.Where("active = ?", *active)
	}

	if searchLabel != "" {
		query = query.Where("label ILIKE ?", "%"+searchLabel+"%") // PostgreSQL
	}

	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if sortBy != "label" && sortBy != "created_at" {
		sortBy = "created_at"
	}

	query = query.Order(sortBy + " " + "asc")

	if err := query.Offset((page - 1) * perPage).Limit(perPage).Find(&urls).Error; err != nil {
		return nil, 0, err
	}

	return urls, int(count), nil
}

func (r *urlRepository) GetActiveURLs(tx *gorm.DB) ([]models.URL, error) {
	var urls []models.URL
	err := tx.Preload("User").Where("active = ?", true).Find(&urls).Error
	if err != nil {
		return nil, err
	}
	return urls, nil
}
