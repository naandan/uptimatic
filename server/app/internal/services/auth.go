package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"uptimatic/internal/email"
	"uptimatic/internal/models"
	"uptimatic/internal/repositories"
	"uptimatic/internal/tasks"
	"uptimatic/internal/utils"

	"github.com/go-redis/redis/v8"
	"github.com/hibiken/asynq"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(ctx context.Context, email, password, appUrl string) (*models.User, *utils.AppError)
	Login(ctx context.Context, email, password string) (string, string, *utils.AppError)
	Logout(ctx context.Context, refreshToken string) *utils.AppError
	Refresh(ctx context.Context, refreshToken string) (string, string, *utils.AppError)
	Profile(ctx context.Context, userId uint) (*models.User, *utils.AppError)
	VerifyEmail(ctx context.Context, token string) *utils.AppError
	ResendVerificationEmail(ctx context.Context, userId uint, appUrl string) (int, *utils.AppError)
	ResendVerificationEmailTTL(ctx context.Context, userId uint) (int, *utils.AppError)
	SendPasswordResetEmail(ctx context.Context, userEmail, appUrl string) *utils.AppError
	ResetPassword(ctx context.Context, token, password string) *utils.AppError
}

type authService struct {
	db          *gorm.DB
	userRepo    repositories.UserRepository
	redis       *redis.Client
	jwtUtil     utils.JWTUtil
	asyncClient *asynq.Client
}

func NewAuthService(db *gorm.DB, userRepo repositories.UserRepository, redis *redis.Client, jwtUtil utils.JWTUtil, asyncClient *asynq.Client) AuthService {
	return &authService{db, userRepo, redis, jwtUtil, asyncClient}
}

func (s *authService) Register(ctx context.Context, userEmail, password string, appUrl string) (*models.User, *utils.AppError) {
	_, err := s.userRepo.FindByEmail(ctx, s.db, userEmail)
	if err == nil {
		return nil, utils.UniqueFieldError("email")
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := &models.User{Email: userEmail, Password: string(hashed)}
	if err := s.userRepo.Create(ctx, s.db, user); err != nil {
		return nil, utils.InternalServerError("Error creating user", err)
	}

	token, err := s.jwtUtil.GenerateEmailVerificationToken(user.ID)
	if err != nil {
		return nil, utils.InternalServerError("Error generating token", err)
	}

	link := fmt.Sprintf("%s/auth/verify?token=%s", appUrl, token)
	iconUrl := fmt.Sprintf("%s/icon.png", appUrl)
	tasks.EnqueueEmail(s.asyncClient, userEmail, "Verify your email - Uptimatic", email.EmailVerify, map[string]any{"LogoURL": iconUrl, "Name": userEmail, "VerificationLink": link})

	if err := s.redis.Set(ctx, utils.GetEmailVerificationTokenKey(userEmail), user.ID, 60*time.Second).Err(); err != nil {
		return nil, utils.InternalServerError("Error generating token", err)
	}

	return user, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (string, string, *utils.AppError) {
	user, err := s.userRepo.FindByEmail(ctx, s.db, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", "", utils.NewAppError(http.StatusUnauthorized, utils.Unauthorized, "email or password incorrect", err)
		}
		return "", "", utils.InternalServerError("Error finding user", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", utils.NewAppError(http.StatusUnauthorized, utils.Unauthorized, "email or password incorrect", err)
	}

	access, refresh, err := s.jwtUtil.GenerateTokens(user.ID, user.Verified)
	if err != nil {
		return "", "", utils.InternalServerError("Error generating tokens", err)
	}

	if err := s.redis.Set(ctx, utils.GetRefreshTokenKey(refresh), user.ID, s.jwtUtil.RefreshTTL).Err(); err != nil {
		return "", "", utils.InternalServerError("Error generating tokens", err)
	}

	return access, refresh, nil
}

func (s *authService) Logout(ctx context.Context, refreshToken string) *utils.AppError {
	err := s.redis.Del(ctx, utils.GetRefreshTokenKey(refreshToken)).Err()
	if err != nil {
		return utils.InternalServerError("Error deleting refresh token", err)
	}
	return nil
}

func (s *authService) Refresh(ctx context.Context, oldRefresh string) (string, string, *utils.AppError) {
	id, err := s.redis.Get(ctx, utils.GetRefreshTokenKey(oldRefresh)).Result()
	if err != nil {
		return "", "", utils.NewAppError(http.StatusUnauthorized, utils.Unauthorized, "Invalid refresh token", err)
	}

	uintID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return "", "", utils.InternalServerError("Error parsing refresh token", err)
	}

	if err := s.redis.Del(ctx, utils.GetRefreshTokenKey(oldRefresh)).Err(); err != nil {
		return "", "", utils.InternalServerError("Error deleting refresh token", err)
	}

	user, err := s.userRepo.FindByID(ctx, s.db, uint(uintID))
	if err != nil {
		return "", "", utils.InternalServerError("Error finding user", err)
	}

	access, refresh, err := s.jwtUtil.GenerateTokens(user.ID, user.Verified)
	if err != nil {
		return "", "", utils.InternalServerError("Error generating tokens", err)
	}

	if err := s.redis.Set(ctx, utils.GetRefreshTokenKey(refresh), user.ID, s.jwtUtil.RefreshTTL).Err(); err != nil {
		return "", "", utils.InternalServerError("Error generating tokens", err)
	}

	return access, refresh, nil
}

func (s *authService) Profile(ctx context.Context, userId uint) (*models.User, *utils.AppError) {
	user, err := s.userRepo.FindByID(ctx, s.db, userId)
	if err != nil {
		return nil, utils.InternalServerError("Error finding user", err)
	}

	return user, nil
}

func (s *authService) VerifyEmail(ctx context.Context, token string) *utils.AppError {
	claims, err := s.jwtUtil.ValidateToken(token)
	if err != nil {
		return utils.NewAppError(http.StatusUnauthorized, utils.InvalidToken, "Invalid or expired token", err)
	}

	id, ok := claims["user_id"].(float64)
	if !ok {
		return utils.NewAppError(http.StatusUnauthorized, utils.InvalidToken, "Invalid or expired token", err)
	}

	userId := uint(id)
	user, err := s.userRepo.FindByID(ctx, s.db, userId)
	if err != nil {
		return utils.InternalServerError("Error finding user", err)
	}

	if user.Verified {
		return utils.ConflictError("Email already verified", err)
	}

	user.Verified = true
	if err := s.userRepo.Update(ctx, s.db, user); err != nil {
		return utils.InternalServerError("Error updating user", err)
	}

	return nil
}

func (s *authService) ResendVerificationEmail(ctx context.Context, userID uint, appUrl string) (int, *utils.AppError) {
	user, err := s.userRepo.FindByID(ctx, s.db, userID)
	if err != nil {
		return 0, utils.InternalServerError("Error finding user", err)
	}
	if user.Verified {
		return 0, utils.ConflictError("Email already verified", err)
	}

	expSec, err := s.redis.TTL(ctx, utils.GetEmailVerificationTokenKey(user.Email)).Result()
	if err == nil && expSec > 0 {
		return int(expSec.Seconds()), nil
	}

	token, err := s.jwtUtil.GenerateEmailVerificationToken(user.ID)
	if err != nil {
		return 0, utils.InternalServerError("Error generating token", err)
	}

	link := fmt.Sprintf("%s/auth/verify?token=%s", appUrl, token)
	iconUrl := fmt.Sprintf("%s/icon.png", appUrl)
	tasks.EnqueueEmail(s.asyncClient, user.Email, "Verify your email - Uptimatic", email.EmailVerify, map[string]any{"LogoURL": iconUrl, "Name": user.Email, "VerificationLink": link})

	if err := s.redis.Set(ctx, utils.GetEmailVerificationTokenKey(user.Email), user.ID, 60*time.Second).Err(); err != nil {
		return 0, utils.InternalServerError("Error generating token", err)
	}

	expSecInt, err := s.redis.TTL(ctx, utils.GetEmailVerificationTokenKey(user.Email)).Result()
	if err != nil {
		return 0, utils.InternalServerError("Error generating token", err)
	}
	return int(expSecInt.Seconds()), nil
}

func (s *authService) ResendVerificationEmailTTL(ctx context.Context, userId uint) (int, *utils.AppError) {
	user, err := s.userRepo.FindByID(ctx, s.db, userId)
	if err != nil {
		return 0, utils.InternalServerError("Error finding user", err)
	}
	if user.Verified {
		return 0, utils.ConflictError("Email already verified", err)
	}

	ttl, err := s.redis.TTL(ctx, utils.GetEmailVerificationTokenKey(user.Email)).Result()
	if err != nil {
		return 0, utils.InternalServerError("Error generating token", err)
	}
	if ttl < 0 {
		return 0, nil
	}
	return int(ttl.Seconds()), nil
}

func (s *authService) SendPasswordResetEmail(ctx context.Context, userEmail, appUrl string) *utils.AppError {
	user, err := s.userRepo.FindByEmail(ctx, s.db, userEmail)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return utils.InternalServerError("Error finding user", err)
	}

	token, err := s.jwtUtil.GeneratePasswordResetToken(user.ID)
	if err != nil {
		return utils.InternalServerError("Error generating token", err)
	}

	link := fmt.Sprintf("%s/auth/reset-password?token=%s", appUrl, token)
	iconUrl := fmt.Sprintf("%s/icon.png", appUrl)
	tasks.EnqueueEmail(s.asyncClient, user.Email, "Reset your password - Uptimatic", email.EmailPasswordReset, map[string]any{"LogoURL": iconUrl, "Name": user.Email, "ResetLink": link})

	return nil
}

func (s *authService) ResetPassword(ctx context.Context, token, password string) *utils.AppError {
	claims, err := s.jwtUtil.ValidateToken(token)
	if err != nil {
		return utils.NewAppError(http.StatusUnauthorized, utils.InvalidToken, "Invalid or expired token", err)
	}

	id, ok := claims["user_id"].(float64)
	if !ok {
		return utils.NewAppError(http.StatusUnauthorized, utils.InvalidToken, "Invalid or expired token", err)
	}

	userId := uint(id)
	user, err := s.userRepo.FindByID(ctx, s.db, userId)
	if err != nil {
		return utils.InternalServerError("Error finding user", err)
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user.Password = string(hashed)
	if err := s.userRepo.Update(ctx, s.db, user); err != nil {
		return utils.InternalServerError("Error updating user", err)
	}

	return nil
}
