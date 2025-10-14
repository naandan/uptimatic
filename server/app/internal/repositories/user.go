package repositories

import (
	"uptimatic/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(tx *gorm.DB, user *models.User) error
	Update(tx *gorm.DB, user *models.User) error
	FindByID(tx *gorm.DB, id uint) (*models.User, error)
	FindByEmail(tx *gorm.DB, email string) (*models.User, error)
}

type userRepository struct{}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (r *userRepository) Create(tx *gorm.DB, user *models.User) error {
	return tx.Create(user).Error
}

func (r *userRepository) Update(tx *gorm.DB, user *models.User) error {
	return tx.Save(user).Error
}

func (r *userRepository) FindByID(tx *gorm.DB, id uint) (*models.User, error) {
	var user models.User
	err := tx.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(tx *gorm.DB, email string) (*models.User, error) {
	var user models.User
	err := tx.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
