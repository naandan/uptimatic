package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"uptimatic/internal/adapters/email"
	"uptimatic/internal/adapters/google"
	"uptimatic/internal/models"
	"uptimatic/internal/tasks"
	"uptimatic/internal/user"
	"uptimatic/internal/utils"

	"github.com/go-redis/redis/v8"
	"github.com/hibiken/asynq"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(ctx context.Context, name, email, password, appUrl string) (*models.User, *utils.AppError)
	Login(ctx context.Context, email, password string) (string, string, *utils.AppError)
	Logout(ctx context.Context, refreshToken string) *utils.AppError
	Refresh(ctx context.Context, refreshToken string) (string, string, *utils.AppError)
	VerifyEmail(ctx context.Context, token string) *utils.AppError
	ResendVerificationEmail(ctx context.Context, userId uint, appUrl string) (int, *utils.AppError)
	ResendVerificationEmailTTL(ctx context.Context, userId uint) (int, *utils.AppError)
	SendPasswordResetEmail(ctx context.Context, userEmail, appUrl string) *utils.AppError
	ResetPassword(ctx context.Context, token, password string) *utils.AppError
	GoogleLogin(ctx context.Context) string
	GoogleCallback(ctx context.Context, code string) (string, string, *utils.AppError)
}

type authService struct {
	db          *gorm.DB
	userRepo    user.UserRepository
	redis       *redis.Client
	jwtUtil     utils.JWTUtil
	asyncClient *asynq.Client
	google      *google.GoogleCLient
}

func NewAuthService(db *gorm.DB, userRepo user.UserRepository, redis *redis.Client, jwtUtil utils.JWTUtil, asyncClient *asynq.Client, google *google.GoogleCLient) AuthService {
	return &authService{db, userRepo, redis, jwtUtil, asyncClient, google}
}

func (s *authService) Register(ctx context.Context, name, userEmail, password, appUrl string) (*models.User, *utils.AppError) {
	utils.Info(ctx, "Register attempt", map[string]any{"email": userEmail})

	_, err := s.userRepo.FindByEmail(ctx, s.db, userEmail)
	if err == nil {
		utils.Warn(ctx, "Register failed: email already exists", map[string]any{"email": userEmail})
		return nil, utils.UniqueFieldError("email")
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := &models.User{Name: name, Email: userEmail, Password: string(hashed)}

	if err := s.userRepo.Create(ctx, s.db, user); err != nil {
		utils.Error(ctx, "Error creating user", map[string]any{"email": userEmail, "err": err.Error()})
		return nil, utils.InternalServerError("Error creating user", err)
	}

	token, err := s.jwtUtil.GenerateEmailVerificationToken(user.ID)
	if err != nil {
		utils.Error(ctx, "Error generating verification token", map[string]any{"user_id": user.ID, "err": err.Error()})
		return nil, utils.InternalServerError("Error generating token", err)
	}

	link := fmt.Sprintf("%s/auth/verify?token=%s", appUrl, token)
	iconUrl := fmt.Sprintf("%s/icon.png", appUrl)
	if err := tasks.EnqueueEmail(s.asyncClient, userEmail, "Verify your email - Uptimatic", email.EmailVerify, map[string]any{
		"LogoURL":          iconUrl,
		"Name":             userEmail,
		"VerificationLink": link,
	}); err != nil {
		utils.Error(ctx, "Error sending verification email", map[string]any{"user_id": user.ID, "email": userEmail, "err": err.Error()})
		return nil, utils.InternalServerError("Error sending verification email", err)
	}

	if err := s.redis.Set(ctx, utils.GetEmailVerificationTokenKey(userEmail), user.ID, 60*time.Second).Err(); err != nil {
		utils.Error(ctx, "Error saving verification token in Redis", map[string]any{"user_id": user.ID, "err": err.Error()})
		return nil, utils.InternalServerError("Error generating token", err)
	}

	utils.Info(ctx, "User registered successfully", map[string]any{"user_id": user.ID, "email": userEmail})
	return user, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (string, string, *utils.AppError) {
	utils.Info(ctx, "Login attempt", map[string]any{"email": email})

	user, err := s.userRepo.FindByEmail(ctx, s.db, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Warn(ctx, "Login failed: user not found", map[string]any{"email": email})
			return "", "", utils.NewAppError(http.StatusUnauthorized, utils.InvalidCredentials, "email or password incorrect", err)
		}
		utils.Error(ctx, "Login DB error", map[string]any{"email": email, "err": err.Error()})
		return "", "", utils.InternalServerError("Error finding user", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		utils.Warn(ctx, "Login failed: wrong password", map[string]any{"email": email})
		return "", "", utils.NewAppError(http.StatusUnauthorized, utils.InvalidCredentials, "email or password incorrect", err)
	}

	access, refresh, err := s.jwtUtil.GenerateTokens(user.ID, user.Verified)
	if err != nil {
		utils.Error(ctx, "Error generating tokens", map[string]any{"user_id": user.ID, "err": err.Error()})
		return "", "", utils.InternalServerError("Error generating tokens", err)
	}

	if err := s.redis.Set(ctx, utils.GetRefreshTokenKey(refresh), user.ID, s.jwtUtil.RefreshTTL).Err(); err != nil {
		utils.Error(ctx, "Error storing refresh token", map[string]any{"user_id": user.ID, "err": err.Error()})
		return "", "", utils.InternalServerError("Error generating tokens", err)
	}

	utils.Info(ctx, "User logged in successfully", map[string]any{"user_id": user.ID, "email": email})
	return access, refresh, nil
}

func (s *authService) Logout(ctx context.Context, refreshToken string) *utils.AppError {
	utils.Info(ctx, "Logout attempt", map[string]any{"refresh_token": refreshToken})

	err := s.redis.Del(ctx, utils.GetRefreshTokenKey(refreshToken)).Err()
	if err != nil {
		utils.Error(ctx, "Error deleting refresh token", map[string]any{"err": err.Error()})
		return utils.InternalServerError("Error deleting refresh token", err)
	}

	utils.Info(ctx, "User logged out successfully", nil)
	return nil
}

func (s *authService) Refresh(ctx context.Context, oldRefresh string) (string, string, *utils.AppError) {
	utils.Info(ctx, "Token refresh attempt", map[string]any{"refresh_token": oldRefresh})

	id, err := s.redis.Get(ctx, utils.GetRefreshTokenKey(oldRefresh)).Result()
	if err != nil {
		utils.Warn(ctx, "Invalid refresh token", map[string]any{"err": err.Error()})
		return "", "", utils.NewAppError(http.StatusUnauthorized, utils.Unauthorized, "Invalid refresh token", err)
	}

	uintID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		utils.Error(ctx, "Error parsing refresh token ID", map[string]any{"err": err.Error()})
		return "", "", utils.InternalServerError("Error parsing refresh token", err)
	}

	if err := s.redis.Del(ctx, utils.GetRefreshTokenKey(oldRefresh)).Err(); err != nil {
		utils.Error(ctx, "Failed to delete old refresh token", map[string]any{"err": err.Error()})
		return "", "", utils.InternalServerError("Error deleting refresh token", err)
	}

	user, err := s.userRepo.FindByID(ctx, s.db, uint(uintID))
	if err != nil {
		utils.Error(ctx, "Error finding user for refresh", map[string]any{"user_id": uintID, "err": err.Error()})
		return "", "", utils.InternalServerError("Error finding user", err)
	}

	access, refresh, err := s.jwtUtil.GenerateTokens(user.ID, user.Verified)
	if err != nil {
		utils.Error(ctx, "Error generating new tokens", map[string]any{"user_id": user.ID, "err": err.Error()})
		return "", "", utils.InternalServerError("Error generating tokens", err)
	}

	if err := s.redis.Set(ctx, utils.GetRefreshTokenKey(refresh), user.ID, s.jwtUtil.RefreshTTL).Err(); err != nil {
		utils.Error(ctx, "Error storing new refresh token", map[string]any{"user_id": user.ID, "err": err.Error()})
		return "", "", utils.InternalServerError("Error generating tokens", err)
	}

	utils.Info(ctx, "Token refreshed successfully", map[string]any{"user_id": user.ID})
	return access, refresh, nil
}

func (s *authService) VerifyEmail(ctx context.Context, token string) *utils.AppError {
	utils.Info(ctx, "Email verification attempt", map[string]any{"token": token})

	claims, err := s.jwtUtil.ValidateToken(token)
	if err != nil {
		utils.Warn(ctx, "Invalid or expired email verification token", map[string]any{"err": err.Error()})
		return utils.NewAppError(http.StatusUnauthorized, utils.InvalidToken, "Invalid or expired token", err)
	}

	id, ok := claims["user_id"].(float64)
	if !ok {
		utils.Error(ctx, "Invalid token structure", nil)
		return utils.NewAppError(http.StatusUnauthorized, utils.InvalidToken, "Invalid or expired token", err)
	}

	userId := uint(id)
	user, err := s.userRepo.FindByID(ctx, s.db, userId)
	if err != nil {
		utils.Error(ctx, "Error finding user during verification", map[string]any{"user_id": userId, "err": err.Error()})
		return utils.InternalServerError("Error finding user", err)
	}

	user.Verified = true
	if err := s.userRepo.Update(ctx, s.db, user); err != nil {
		utils.Error(ctx, "Error updating user verification", map[string]any{"user_id": userId, "err": err.Error()})
		return utils.InternalServerError("Error updating user", err)
	}

	utils.Info(ctx, "Email verified successfully", map[string]any{"user_id": userId, "email": user.Email})
	return nil
}

func (s *authService) ResendVerificationEmail(ctx context.Context, userID uint, appUrl string) (int, *utils.AppError) {
	utils.Info(ctx, "Resend verification email attempt", map[string]any{"user_id": userID})

	user, err := s.userRepo.FindByID(ctx, s.db, userID)
	if err != nil {
		utils.Error(ctx, "Error finding user for resend verification", map[string]any{"user_id": userID, "err": err.Error()})
		return 0, utils.InternalServerError("Error finding user", err)
	}

	ttl, err := s.redis.TTL(ctx, utils.GetEmailVerificationTokenKey(user.Email)).Result()
	if err == nil && ttl > 0 {
		utils.Info(ctx, "Resend verification email throttled", map[string]any{"user_id": userID, "ttl": ttl.Seconds()})
		return int(ttl.Seconds()), nil
	}

	token, err := s.jwtUtil.GenerateEmailVerificationToken(user.ID)
	if err != nil {
		utils.Error(ctx, "Error generating new verification token", map[string]any{"user_id": userID, "err": err.Error()})
		return 0, utils.InternalServerError("Error generating token", err)
	}

	link := fmt.Sprintf("%s/auth/verify?token=%s", appUrl, token)
	iconUrl := fmt.Sprintf("%s/icon.png", appUrl)
	if err := tasks.EnqueueEmail(s.asyncClient, user.Email, "Verify your email - Uptimatic", email.EmailVerify, map[string]any{
		"LogoURL":          iconUrl,
		"Name":             user.Email,
		"VerificationLink": link,
	}); err != nil {
		utils.Error(ctx, "Error sending verification email", map[string]any{"user_id": userID, "err": err.Error()})
		return 0, utils.InternalServerError("Error sending verification email", err)
	}

	if err := s.redis.Set(ctx, utils.GetEmailVerificationTokenKey(user.Email), user.ID, 60*time.Second).Err(); err != nil {
		utils.Error(ctx, "Error storing new verification token", map[string]any{"user_id": userID, "err": err.Error()})
		return 0, utils.InternalServerError("Error generating token", err)
	}

	expSecInt, _ := s.redis.TTL(ctx, utils.GetEmailVerificationTokenKey(user.Email)).Result()
	utils.Info(ctx, "Verification email resent successfully", map[string]any{"user_id": userID, "ttl": expSecInt.Seconds()})
	return int(expSecInt.Seconds()), nil
}

func (s *authService) ResendVerificationEmailTTL(ctx context.Context, userId uint) (int, *utils.AppError) {
	utils.Info(ctx, "Checking resend verification TTL", map[string]any{"user_id": userId})

	user, err := s.userRepo.FindByID(ctx, s.db, userId)
	if err != nil {
		utils.Error(ctx, "Error finding user for TTL check", map[string]any{"user_id": userId, "err": err.Error()})
		return 0, utils.InternalServerError("Error finding user", err)
	}

	ttl, err := s.redis.TTL(ctx, utils.GetEmailVerificationTokenKey(user.Email)).Result()
	if err != nil {
		utils.Error(ctx, "Error getting TTL from Redis", map[string]any{"user_id": userId, "err": err.Error()})
		return 0, utils.InternalServerError("Error generating token", err)
	}
	if ttl < 0 {
		utils.Warn(ctx, "No verification token found in Redis", map[string]any{"user_id": userId})
		return 0, nil
	}

	utils.Info(ctx, "Verification token TTL retrieved", map[string]any{"user_id": userId, "ttl": ttl.Seconds()})
	return int(ttl.Seconds()), nil
}

func (s *authService) SendPasswordResetEmail(ctx context.Context, userEmail, appUrl string) *utils.AppError {
	utils.Info(ctx, "Password reset email request", map[string]any{"email": userEmail})

	user, err := s.userRepo.FindByEmail(ctx, s.db, userEmail)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Warn(ctx, "Password reset requested for non-existent email", map[string]any{"email": userEmail})
			return nil
		}
		utils.Error(ctx, "Error finding user for password reset", map[string]any{"email": userEmail, "err": err.Error()})
		return utils.InternalServerError("Error finding user", err)
	}

	token, err := s.jwtUtil.GeneratePasswordResetToken(user.ID)
	if err != nil {
		utils.Error(ctx, "Error generating password reset token", map[string]any{"user_id": user.ID, "err": err.Error()})
		return utils.InternalServerError("Error generating token", err)
	}

	link := fmt.Sprintf("%s/auth/reset-password?token=%s", appUrl, token)
	iconUrl := fmt.Sprintf("%s/icon.png", appUrl)
	if err := tasks.EnqueueEmail(s.asyncClient, user.Email, "Reset your password - Uptimatic", email.EmailPasswordReset, map[string]any{
		"LogoURL":   iconUrl,
		"Name":      user.Email,
		"ResetLink": link,
	}); err != nil {
		utils.Error(ctx, "Error sending password reset email", map[string]any{"user_id": user.ID, "err": err.Error()})
		return utils.InternalServerError("Error sending password reset email", err)
	}

	utils.Info(ctx, "Password reset email queued", map[string]any{"user_id": user.ID, "email": user.Email})
	return nil
}

func (s *authService) ResetPassword(ctx context.Context, token, password string) *utils.AppError {
	utils.Info(ctx, "Password reset attempt", nil)

	claims, err := s.jwtUtil.ValidateToken(token)
	if err != nil {
		utils.Warn(ctx, "Invalid or expired password reset token", map[string]any{"err": err.Error()})
		return utils.NewAppError(http.StatusUnauthorized, utils.InvalidToken, "Invalid or expired token", err)
	}

	id, ok := claims["user_id"].(float64)
	if !ok {
		utils.Error(ctx, "Invalid token structure", nil)
		return utils.NewAppError(http.StatusUnauthorized, utils.InvalidToken, "Invalid or expired token", err)
	}

	userId := uint(id)
	user, err := s.userRepo.FindByID(ctx, s.db, userId)
	if err != nil {
		utils.Error(ctx, "Error finding user during password reset", map[string]any{"user_id": userId, "err": err.Error()})
		return utils.InternalServerError("Error finding user", err)
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user.Password = string(hashed)

	if err := s.userRepo.Update(ctx, s.db, user); err != nil {
		utils.Error(ctx, "Error updating password", map[string]any{"user_id": userId, "err": err.Error()})
		return utils.InternalServerError("Error updating user", err)
	}

	utils.Info(ctx, "Password reset successfully", map[string]any{"user_id": userId})
	return nil
}

func (s *authService) GoogleLogin(ctx context.Context) string {
	url := s.google.AuthCodeURL("randomstate")
	return url
}

func (s *authService) GoogleCallback(ctx context.Context, code string) (string, string, *utils.AppError) {
	token, err := s.google.Exchange(ctx, code)
	if err != nil {
		utils.Error(ctx, "Token exchange failed", map[string]any{"err": err.Error()})
		return "", "", utils.InternalServerError("Token exchange failed", err)
	}

	client := s.google.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		utils.Error(ctx, "Failed to get user info", map[string]any{"err": err.Error()})
		return "", "", utils.InternalServerError("Failed to get user info", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			utils.Error(ctx, "Failed to close response body", map[string]any{"err": err})
		}
	}()

	var oauthUser map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&oauthUser); err != nil {
		utils.Error(ctx, "Failed to decode user info", map[string]any{"err": err.Error()})
		return "", "", utils.InternalServerError("Failed to decode user info", err)
	}

	email := oauthUser["email"].(string)
	user, err := s.userRepo.FindByEmail(ctx, s.db, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			user = &models.User{Name: oauthUser["name"].(string), Email: email, Verified: true}
			if err := s.userRepo.Create(ctx, s.db, user); err != nil {
				utils.Error(ctx, "Failed to create user", map[string]any{"err": err.Error()})
				return "", "", utils.InternalServerError("Failed to create user", err)
			}
		} else {
			utils.Error(ctx, "Failed to find user", map[string]any{"err": err.Error()})
			return "", "", utils.InternalServerError("Failed to find user", err)
		}
	}

	access, refresh, err := s.jwtUtil.GenerateTokens(user.ID, user.Verified)
	if err != nil {
		utils.Error(ctx, "Failed to generate token", map[string]any{"err": err.Error()})
		return "", "", utils.InternalServerError("Failed to generate token", err)
	}

	if err := s.redis.Set(ctx, utils.GetRefreshTokenKey(refresh), user.ID, s.jwtUtil.RefreshTTL).Err(); err != nil {
		utils.Error(ctx, "Failed to store refresh token", map[string]any{"err": err.Error()})
		return "", "", utils.InternalServerError("Failed to store refresh token", err)
	}

	return access, refresh, nil
}
