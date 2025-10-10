package services

import (
	"context"
	"errors"
	"strconv"
	"uptimatic/internal/models"
	"uptimatic/internal/repositories"
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
	VerifyEmail(token string) error
	ResendVerificationEmail(email string, appUrl string) error
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

	// token, err := s.jwtUtil.GenerateEmailVerificationToken(user.ID)
	// if err != nil {
	// 	return nil, err
	// }

	// link := fmt.Sprintf("%s/api/v1/auth/verify-email?token=%s", appUrl, token)
	// tasks.EnqueueEmail(s.asyncClient, userEmail, "Verify your email - Cartel", email.EmailVerify, map[string]any{"Name": name, "VerificationLink": link})
	return user, nil
}

func (s *authService) Login(email, password string) (string, string, error) {
	user, err := s.userRepo.FindByEmail(s.db, email)
	if err != nil {
		return "", "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", err
	}

	access, refresh, err := s.jwtUtil.GenerateTokens(user.ID)
	if err != nil {
		return "", "", err
	}

	// rt := &models.RefreshToken{
	// 	UserID:    user.ID,
	// 	Token:     refresh,
	// 	ExpiresAt: time.Now().Add(s.jwtUtil.RefreshTTL),
	// }
	// if err := s.refreshRepo.Save(rt); err != nil {
	// 	return "", "", err
	// }

	if err := s.redis.Set(context.Background(), utils.GetRefreshTokenKey(refresh), user.ID, s.jwtUtil.RefreshTTL).Err(); err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

func (s *authService) Logout(refreshToken string) error {
	return s.redis.Del(context.Background(), utils.GetRefreshTokenKey(refreshToken)).Err()
}

func (s *authService) Refresh(oldRefresh string) (string, string, error) {
	// rt, err := s.refreshRepo.FindValid(oldRefresh)
	// if err != nil {
	// 	return "", "", err
	// }

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

	access, refresh, err := s.jwtUtil.GenerateTokens(user.ID)
	if err != nil {
		return "", "", err
	}

	// newRT := &models.RefreshToken{
	// 	UserID:    rt.UserID,
	// 	Token:     refresh,
	// 	ExpiresAt: time.Now().Add(s.jwtUtil.RefreshTTL),
	// }
	// if err := s.refreshRepo.Save(newRT); err != nil {
	// 	return "", "", err
	// }

	if err := s.redis.Set(context.Background(), utils.GetRefreshTokenKey(refresh), user.ID, s.jwtUtil.RefreshTTL).Err(); err != nil {
		return "", "", err
	}

	return access, refresh, nil
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

func (s *authService) ResendVerificationEmail(userEmail string, appUrl string) error {
	user, err := s.userRepo.FindByEmail(s.db, userEmail)
	if err != nil {
		return err
	}
	if user.Verified {
		return errors.New("email already verified")
	}

	// token, err := s.jwtUtil.GenerateEmailVerificationToken(user.ID)
	// if err != nil {
	// 	return err
	// }

	// link := fmt.Sprintf("%s/api/v1/auth/verify-email?token=%s", appUrl, token)
	// tasks.EnqueueEmail(s.asyncClient, userEmail, "Verify your email - Cartel", email.EmailVerify, map[string]any{"Name": user.Name, "VerificationLink": link})
	return nil
}
