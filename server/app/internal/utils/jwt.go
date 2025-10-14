package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTUtil struct {
	Secret     string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

func NewJWTUtil(secret string, accessTTL, refreshTTL time.Duration) JWTUtil {
	return JWTUtil{
		Secret:     secret,
		AccessTTL:  accessTTL,
		RefreshTTL: refreshTTL,
	}
}

func (j *JWTUtil) GenerateTokens(userID uint, verified bool) (string, string, error) {
	now := time.Now()

	accessClaims := jwt.MapClaims{
		"user_id":  userID,
		"verified": verified,
		"exp":      now.Add(j.AccessTTL).Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	access, err := accessToken.SignedString([]byte(j.Secret))
	if err != nil {
		return "", "", err
	}

	refreshClaims := jwt.MapClaims{
		"user_id":  userID,
		"verified": verified,
		"exp":      now.Add(j.RefreshTTL).Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refresh, err := refreshToken.SignedString([]byte(j.Secret))
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

func (j *JWTUtil) ParseToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		return []byte(j.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

func (j *JWTUtil) ValidateToken(tokenStr string) (jwt.MapClaims, error) {
	claims, err := j.ParseToken(tokenStr)
	if err != nil {
		return nil, err
	}
	if exp, ok := claims["exp"].(float64); ok && exp < float64(time.Now().Unix()) {
		return nil, jwt.ErrTokenExpired
	}
	return claims, nil
}

func (j *JWTUtil) GenerateEmailVerificationToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(2 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.Secret))
}

func (j *JWTUtil) GeneratePasswordResetToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(2 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.Secret))
}
