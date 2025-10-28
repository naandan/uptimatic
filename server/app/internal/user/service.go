package user

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"time"
	"uptimatic/internal/adapters/email"
	"uptimatic/internal/adapters/minio"
	"uptimatic/internal/models"
	"uptimatic/internal/tasks"
	"uptimatic/internal/utils"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	Update(ctx context.Context, userId uint, name, userEmail, appUrl, oldRefresh string) (*models.User, map[string]any, *utils.AppError)
	GetUser(ctx context.Context, userId uint) (*models.User, *utils.AppError)
	ChangePassword(ctx context.Context, userId uint, oldPassword, newPassword string) *utils.AppError
	GetPresignedUrl(ctx context.Context, fileName string, contentType string) (string, string, *utils.AppError)
	UpdateFoto(ctx context.Context, userId uint, fileName string) (string, *utils.AppError)
}

type userService struct {
	db          *gorm.DB
	userRepo    UserRepository
	minio       *minio.MinioUtil
	redis       *redis.Client
	jwtUtil     utils.JWTUtil
	asyncClient *asynq.Client
}

func NewUserService(db *gorm.DB, userRepo UserRepository, minio *minio.MinioUtil, redis *redis.Client, jwtUtil utils.JWTUtil, asyncClient *asynq.Client) UserService {
	return &userService{db, userRepo, minio, redis, jwtUtil, asyncClient}
}

func (s *userService) Update(ctx context.Context, userId uint, name, userEmail, appUrl, oldRefresh string) (*models.User, map[string]any, *utils.AppError) {
	utils.Info(ctx, "Updating user profile", map[string]any{"name": name, "email": userEmail})

	user, err := s.userRepo.FindByID(ctx, s.db, userId)
	if err != nil {
		utils.Error(ctx, "Error updating user profile", map[string]any{"user_id": userId, "err": err.Error()})
		return nil, nil, utils.InternalServerError("Error finding user", err)
	}

	if userEmail == user.Email {
		user.Name = name
		if err := s.userRepo.Update(ctx, s.db, user); err != nil {
			utils.Error(ctx, "Error updating user profile", map[string]any{"user_id": user.ID, "err": err.Error()})
			return nil, nil, utils.InternalServerError("Error updating user", err)
		}

		if user.Profile != "" {
			url := s.minio.GetPublicURL(ctx, user.Profile)
			user.Profile = url
		}

		utils.Info(ctx, "User profile updated successfully", map[string]any{"user_id": user.ID})
		return user, nil, nil
	}

	user.Email = userEmail
	user.Verified = false

	token, err := s.jwtUtil.GenerateEmailVerificationToken(user.ID)
	if err != nil {
		utils.Error(ctx, "Error generating verification token", map[string]any{"user_id": user.ID, "err": err.Error()})
		return nil, nil, utils.InternalServerError("Error generating token", err)
	}

	link := fmt.Sprintf("%s/auth/verify?token=%s", appUrl, token)
	iconUrl := fmt.Sprintf("%s/icon.png", appUrl)
	if err := tasks.EnqueueEmail(s.asyncClient, user.Email, "Verify your email - Uptimatic", email.EmailVerify, map[string]any{
		"LogoURL":          iconUrl,
		"Name":             user.Email,
		"VerificationLink": link,
	}); err != nil {
		utils.Error(ctx, "Error sending verification email", map[string]any{"user_id": user.ID, "email": user.Email, "err": err.Error()})
		return nil, nil, utils.InternalServerError("Error sending verification email", err)
	}

	if err := s.redis.Set(ctx, utils.GetEmailVerificationTokenKey(user.Email), user.ID, 60*time.Second).Err(); err != nil {
		utils.Error(ctx, "Error saving verification token in Redis", map[string]any{"user_id": user.ID, "err": err.Error()})
		return nil, nil, utils.InternalServerError("Error generating token", err)
	}

	user.Name = name
	if err := s.userRepo.Update(ctx, s.db, user); err != nil {
		utils.Error(ctx, "Error updating user profile", map[string]any{"user_id": user.ID, "err": err.Error()})
		return nil, nil, utils.InternalServerError("Error updating user", err)
	}

	if err := s.redis.Del(ctx, utils.GetRefreshTokenKey(oldRefresh)).Err(); err != nil {
		utils.Error(ctx, "Error deleting refresh token", map[string]any{"user_id": user.ID, "err": err.Error()})
		return nil, nil, utils.InternalServerError("Error deleting refresh token", err)
	}

	access, refresh, err := s.jwtUtil.GenerateTokens(user.ID, user.Verified)
	if err != nil {
		utils.Error(ctx, "Error generating tokens", map[string]any{"user_id": user.ID, "err": err.Error()})
		return nil, nil, utils.InternalServerError("Error generating tokens", err)
	}

	if err := s.redis.Set(ctx, utils.GetRefreshTokenKey(refresh), user.ID, s.jwtUtil.RefreshTTL).Err(); err != nil {
		utils.Error(ctx, "Error storing refresh token", map[string]any{"user_id": user.ID, "err": err.Error()})
		return nil, nil, utils.InternalServerError("Error generating tokens", err)
	}

	if user.Profile != "" {
		url := s.minio.GetPublicURL(ctx, user.Profile)
		user.Profile = url
	}

	utils.Info(ctx, "User profile updated successfully", map[string]any{"user_id": user.ID})
	changed := map[string]any{
		"is_email_changed": true,
		"access_token":     access,
		"refresh_token":    refresh,
	}
	return user, changed, nil
}

func (s *userService) GetUser(ctx context.Context, userId uint) (*models.User, *utils.AppError) {
	utils.Info(ctx, "Fetching user profile", map[string]any{"user_id": userId})

	user, err := s.userRepo.FindByID(ctx, s.db, userId)
	if err != nil {
		utils.Error(ctx, "Error retrieving user profile", map[string]any{"user_id": userId, "err": err.Error()})
		return nil, utils.InternalServerError("Error finding user", err)
	}

	if user.Profile != "" {
		url := s.minio.GetPublicURL(ctx, user.Profile)
		user.Profile = url
	}

	utils.Info(ctx, "User profile fetched successfully", map[string]any{"user_id": userId})
	return user, nil
}

func (s *userService) ChangePassword(ctx context.Context, userId uint, oldPassword, newPassword string) *utils.AppError {
	utils.Info(ctx, "Changing user password", map[string]any{"user_id": userId})

	user, err := s.userRepo.FindByID(ctx, s.db, userId)
	if err != nil {
		return utils.InternalServerError("Error finding user", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return utils.NewAppError(http.StatusUnauthorized, utils.InvalidCredentials, "Invalid old password", err)
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	user.Password = string(hashed)

	if err := s.userRepo.Update(ctx, s.db, user); err != nil {
		return utils.InternalServerError("Error updating password", err)
	}

	utils.Info(ctx, "Password updated successfully", map[string]any{"user_id": userId})
	return nil
}

func (s *userService) GetPresignedUrl(ctx context.Context, fileName, contentType string) (string, string, *utils.AppError) {
	fileName = uuid.New().String() + filepath.Ext(fileName)
	utils.Info(ctx, "Generating presigned URL for upload", map[string]any{"file": fileName})

	uploadURL, err := s.minio.PutPresignedURL(ctx, fileName, contentType)
	if err != nil {
		return "", "", utils.InternalServerError("Error generating presigned URL", err)
	}

	return uploadURL, fileName, nil
}

func (s *userService) UpdateFoto(ctx context.Context, userId uint, fileName string) (string, *utils.AppError) {
	utils.Info(ctx, "Updating user photo", map[string]any{"user_id": userId})

	user, err := s.userRepo.FindByID(ctx, s.db, userId)
	if err != nil {
		return "", utils.InternalServerError("Error finding user", err)
	}

	if user.Profile != "" {
		if err := s.minio.DeleteFile(ctx, user.Profile); err != nil {
			return "", utils.InternalServerError("Error deleting old photo", err)
		}
	}

	user.Profile = fileName
	if err := s.userRepo.Update(ctx, s.db, user); err != nil {
		return "", utils.InternalServerError("Error updating user photo", err)
	}

	url := s.minio.GetPublicURL(ctx, fileName)

	utils.Info(ctx, "User photo updated successfully", map[string]any{"user_id": userId})
	return url, nil
}
