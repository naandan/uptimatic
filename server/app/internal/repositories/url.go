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
	ListByUserID(tx *gorm.DB, userID uint, page, perPage int) ([]models.URL, int, error)
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

func (r *urlRepository) ListByUserID(tx *gorm.DB, userID uint, page, perPage int) ([]models.URL, int, error) {
	var urls []models.URL
	var count int64
	err := tx.Where("user_id = ?", userID).Offset((page - 1) * perPage).Limit(perPage).Find(&urls).Error
	if err != nil {
		return nil, 0, err
	}
	err = tx.Model(&models.URL{}).Where("user_id = ?", userID).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	return urls, int(count), nil
}

func (r *urlRepository) GetActiveURLs(tx *gorm.DB) ([]models.URL, error) {
	var urls []models.URL
	err := tx.Where("active = ?", true).Find(&urls).Error
	if err != nil {
		return nil, err
	}
	return urls, nil
}
