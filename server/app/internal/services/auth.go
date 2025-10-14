package services

import (
	"context"
	"errors"
	"fmt"
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
	Register(email, password, appUrl string) (*models.User, error)
	Login(email, password string) (string, string, error)
	Logout(refreshToken string) error
	Refresh(refreshToken string) (string, string, error)
	Profile(userId uint) (*models.User, error)
	VerifyEmail(token string) error
	ResendVerificationEmail(userId uint, appUrl string) (int, error)
	ResendVerificationEmailTTL(userId uint) (int, error)
	SendPasswordResetEmail(userEmail, appUrl string) error
	ResetPassword(token, password string) error
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

func (s *authService) Register(userEmail, password string, appUrl string) (*models.User, error) {
	_, err := s.userRepo.FindByEmail(s.db, userEmail)
	if err == nil {
		return nil, errors.New("email already used")
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := &models.User{Email: userEmail, Password: string(hashed)}
	if err := s.userRepo.Create(s.db, user); err != nil {
		return nil, err
	}

	token, err := s.jwtUtil.GenerateEmailVerificationToken(user.ID)
	if err != nil {
		return nil, err
	}

	link := fmt.Sprintf("%s/api/v1/auth/verify?token=%s", appUrl, token)
	tasks.EnqueueEmail(s.asyncClient, userEmail, "Verify your email - Uptimatic", email.EmailVerify, map[string]any{"Name": userEmail, "VerificationLink": link})

	if err := s.redis.Set(context.Background(), utils.GetEmailVerificationTokenKey(userEmail), user.ID, 60*time.Second).Err(); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(email, password string) (string, string, error) {
	user, err := s.userRepo.FindByEmail(s.db, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", "", errors.New("email or password incorrect")
		}
		return "", "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", errors.New("email or password incorrect")
	}

	access, refresh, err := s.jwtUtil.GenerateTokens(user.ID, user.Verified)
	if err != nil {
		return "", "", err
	}

	if err := s.redis.Set(context.Background(), utils.GetRefreshTokenKey(refresh), user.ID, s.jwtUtil.RefreshTTL).Err(); err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

func (s *authService) Logout(refreshToken string) error {
	return s.redis.Del(context.Background(), utils.GetRefreshTokenKey(refreshToken)).Err()
}

func (s *authService) Refresh(oldRefresh string) (string, string, error) {
	id, err := s.redis.Get(context.Background(), utils.GetRefreshTokenKey(oldRefresh)).Result()
	if err != nil {
		return "", "", errors.New("unauthorized")
	}

	uintID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return "", "", err
	}

	if err := s.redis.Del(context.Background(), utils.GetRefreshTokenKey(oldRefresh)).Err(); err != nil {
		return "", "", err
	}

	user, err := s.userRepo.FindByID(s.db, uint(uintID))
	if err != nil {
		return "", "", err
	}

	access, refresh, err := s.jwtUtil.GenerateTokens(user.ID, user.Verified)
	if err != nil {
		return "", "", err
	}

	if err := s.redis.Set(context.Background(), utils.GetRefreshTokenKey(refresh), user.ID, s.jwtUtil.RefreshTTL).Err(); err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

func (s *authService) Profile(userId uint) (*models.User, error) {
	return s.userRepo.FindByID(s.db, userId)
}

func (s *authService) VerifyEmail(token string) error {
	claims, err := s.jwtUtil.ValidateToken(token)
	if err != nil {
		return err
	}

	id, ok := claims["user_id"].(float64)
	if !ok {
		return errors.New("invalid token")
	}

	userId := uint(id)
	user, err := s.userRepo.FindByID(s.db, userId)
	if err != nil {
		return err
	}

	if user.Verified {
		return errors.New("email already verified")
	}

	user.Verified = true
	if err := s.userRepo.Update(s.db, user); err != nil {
		return err
	}

	return nil
}

func (s *authService) ResendVerificationEmail(userID uint, appUrl string) (int, error) {
	user, err := s.userRepo.FindByID(s.db, userID)
	if err != nil {
		return 0, err
	}
	if user.Verified {
		return 0, errors.New("email already verified")
	}

	expSec, err := s.redis.TTL(context.Background(), utils.GetEmailVerificationTokenKey(user.Email)).Result()
	if err == nil && expSec > 0 {
		return int(expSec.Seconds()), nil
	}

	token, err := s.jwtUtil.GenerateEmailVerificationToken(user.ID)
	if err != nil {
		return 0, err
	}

	link := fmt.Sprintf("%s/auth/verify?token=%s", appUrl, token)
	tasks.EnqueueEmail(s.asyncClient, user.Email, "Verify your email - Uptimatic", email.EmailVerify, map[string]any{"Name": user.Email, "VerificationLink": link})

	if err := s.redis.Set(context.Background(), utils.GetEmailVerificationTokenKey(user.Email), user.ID, 60*time.Second).Err(); err != nil {
		return 0, err
	}

	expSecInt, err := s.redis.TTL(context.Background(), utils.GetEmailVerificationTokenKey(user.Email)).Result()
	if err != nil {
		return 0, err
	}
	return int(expSecInt.Seconds()), nil
}

func (s *authService) ResendVerificationEmailTTL(userId uint) (int, error) {
	user, err := s.userRepo.FindByID(s.db, userId)
	if err != nil {
		return 0, err
	}
	if user.Verified {
		return 0, errors.New("email already verified")
	}

	ttl, err := s.redis.TTL(context.Background(), utils.GetEmailVerificationTokenKey(user.Email)).Result()
	if err != nil {
		return 0, err
	}
	if ttl < 0 {
		return 0, nil
	}
	return int(ttl.Seconds()), nil
}

func (s *authService) SendPasswordResetEmail(userEmail, appUrl string) error {
	user, err := s.userRepo.FindByEmail(s.db, userEmail)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	token, err := s.jwtUtil.GeneratePasswordResetToken(user.ID)
	if err != nil {
		return err
	}

	link := fmt.Sprintf("%s/auth/reset-password?token=%s", appUrl, token)
	tasks.EnqueueEmail(s.asyncClient, user.Email, "Reset your password - Uptimatic", email.EmailPasswordReset, map[string]any{"Name": user.Email, "ResetLink": link})

	return nil
}

func (s *authService) ResetPassword(token, password string) error {
	claims, err := s.jwtUtil.ValidateToken(token)
	if err != nil {
		return err
	}

	id, ok := claims["user_id"].(float64)
	if !ok {
		return errors.New("invalid token")
	}

	userId := uint(id)
	user, err := s.userRepo.FindByID(s.db, userId)
	if err != nil {
		return err
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user.Password = string(hashed)
	if err := s.userRepo.Update(s.db, user); err != nil {
		return err
	}

	return nil
}
