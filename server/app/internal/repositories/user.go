package repositories

import (
	"context"
	"uptimatic/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, tx *gorm.DB, user *models.User) error
	Update(ctx context.Context, tx *gorm.DB, user *models.User) error
	FindByID(ctx context.Context, tx *gorm.DB, id uint) (*models.User, error)
	FindByEmail(ctx context.Context, tx *gorm.DB, email string) (*models.User, error)
}

type userRepository struct{}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (r *userRepository) Create(ctx context.Context, tx *gorm.DB, user *models.User) error {
	return tx.WithContext(ctx).Create(user).Error
}

func (r *userRepository) Update(ctx context.Context, tx *gorm.DB, user *models.User) error {
	return tx.WithContext(ctx).Save(user).Error
}

func (r *userRepository) FindByID(ctx context.Context, tx *gorm.DB, id uint) (*models.User, error) {
	var user models.User
	err := tx.WithContext(ctx).First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, tx *gorm.DB, email string) (*models.User, error) {
	var user models.User
	err := tx.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
