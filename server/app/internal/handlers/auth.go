package handlers

import (
	"fmt"
	"net/http"
	"uptimatic/internal/config"
	"uptimatic/internal/schema"
	"uptimatic/internal/services"
	"uptimatic/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthHandler interface {
	RegisterHandler(c *gin.Context)
	LoginHandler(c *gin.Context)
	LogoutHandler(c *gin.Context)
	RefreshHandler(c *gin.Context)
	ProfileHandler(c *gin.Context)
	VerifyHandler(c *gin.Context)
	ResendVerificationHandler(c *gin.Context)
	ResendVerificationEmailTTLHandler(c *gin.Context)
	SendPasswordResetEmailHandler(c *gin.Context)
	ResetPasswordHandler(c *gin.Context)
	GoogleLoginHandler(c *gin.Context)
	GoogleCallbackHandler(c *gin.Context)
}

type authHandler struct {
	authService services.AuthService
	validate    *validator.Validate
	cfg         *config.Config
}

func NewAuthHandler(authService services.AuthService, validate *validator.Validate, cfg *config.Config) AuthHandler {
	return &authHandler{authService, validate, cfg}
}

func (h *authHandler) RegisterHandler(c *gin.Context) {
	var req schema.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, utils.NewAppError(http.StatusBadRequest, utils.ValidationError, "Invalid JSON payload", err))
		return
	}

	if err := h.validate.Struct(req); err != nil {
		utils.BindErrorResponse(c, utils.NewAppError(http.StatusBadRequest, utils.ValidationError, err.Error(), err))
		return
	}

	user, err := h.authService.Register(c.Request.Context(), req.Email, req.Password, fmt.Sprintf("%s://%s", h.cfg.AppScheme, h.cfg.AppDomain))
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, user)
}

func (h *authHandler) LoginHandler(c *gin.Context) {
	var req schema.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, utils.NewAppError(http.StatusBadRequest, utils.ValidationError, "Invalid JSON payload", err))
		return
	}

	if err := h.validate.Struct(req); err != nil {
		utils.BindErrorResponse(c, utils.NewAppError(http.StatusBadRequest, utils.ValidationError, err.Error(), err))
		return
	}

	access, refresh, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	c.SetCookie("access_token", access, int(h.cfg.AuthAccessTokenExpiration), "/", h.cfg.AppDomain, true, true)
	c.SetCookie("refresh_token", refresh, int(h.cfg.AuthRefreshTokenExpiration), "/", h.cfg.AppDomain, true, true)
	utils.SuccessResponse(c, gin.H{
		"access_token":  access,
		"refresh_token": refresh,
	})
}

func (h *authHandler) LogoutHandler(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		utils.ErrorResponse(c, utils.NewAppError(http.StatusUnauthorized, utils.Unauthorized, "Missing refresh token", nil))
		return
	}

	if err := h.authService.Logout(c.Request.Context(), refreshToken); err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	c.SetCookie("access_token", "", -1, "/", h.cfg.AppDomain, true, true)
	c.SetCookie("refresh_token", "", -1, "/", h.cfg.AppDomain, true, true)
	utils.SuccessResponse(c, nil)
}

func (h *authHandler) RefreshHandler(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.SetCookie("access_token", "", -1, "/", h.cfg.AppDomain, true, true)
		c.SetCookie("refresh_token", "", -1, "/", h.cfg.AppDomain, true, true)
		utils.ErrorResponse(c, utils.NewAppError(http.StatusUnauthorized, utils.InvalidCredentials, "Missing refresh token", nil))
		return
	}

	access, refresh, errRef := h.authService.Refresh(c.Request.Context(), refreshToken)
	if errRef != nil {
		c.SetCookie("access_token", "", -1, "/", h.cfg.AppDomain, true, true)
		c.SetCookie("refresh_token", "", -1, "/", h.cfg.AppDomain, true, true)
		utils.ErrorResponse(c, errRef)
		return
	}

	c.SetCookie("access_token", access, int(h.cfg.AuthAccessTokenExpiration), "/", h.cfg.AppDomain, true, true)
	c.SetCookie("refresh_token", refresh, int(h.cfg.AuthRefreshTokenExpiration), "/", h.cfg.AppDomain, true, true)
	utils.SuccessResponse(c, gin.H{
		"access_token":  access,
		"refresh_token": refresh,
	})
}

func (h *authHandler) ProfileHandler(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, utils.NewAppError(http.StatusUnauthorized, utils.Unauthorized, "User not authenticated", nil))
		return
	}

	userIdUint, ok := userId.(uint)
	if !ok {
		utils.ErrorResponse(c, utils.NewAppError(http.StatusUnauthorized, utils.Unauthorized, "User not authenticated", nil))
		return
	}

	user, err := h.authService.Profile(c.Request.Context(), userIdUint)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, user)
}

func (h *authHandler) VerifyHandler(c *gin.Context) {
	token := c.Query("token")
	if err := h.authService.VerifyEmail(c.Request.Context(), token); err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, nil)
}

func (h *authHandler) ResendVerificationHandler(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		utils.ErrorResponse(c, utils.NewAppError(http.StatusUnauthorized, utils.Unauthorized, "User not authenticated", nil))
		return
	}

	ttl, err := h.authService.ResendVerificationEmail(c.Request.Context(), userID, fmt.Sprintf("%s://%s", h.cfg.AppScheme, h.cfg.AppDomain))
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, gin.H{"ttl": ttl})
}

func (h *authHandler) ResendVerificationEmailTTLHandler(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		utils.ErrorResponse(c, utils.NewAppError(http.StatusUnauthorized, utils.Unauthorized, "User not authenticated", nil))
		return
	}

	ttl, err := h.authService.ResendVerificationEmailTTL(c.Request.Context(), userID)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, gin.H{"ttl": ttl})
}

func (h *authHandler) SendPasswordResetEmailHandler(c *gin.Context) {
	var req schema.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, utils.NewAppError(http.StatusBadRequest, utils.ValidationError, "Invalid JSON payload", err))
		return
	}
	if err := h.validate.Struct(req); err != nil {
		utils.BindErrorResponse(c, utils.NewAppError(http.StatusBadRequest, utils.ValidationError, err.Error(), err))
		return
	}
	if err := h.authService.SendPasswordResetEmail(c.Request.Context(), req.Email, fmt.Sprintf("%s://%s", h.cfg.AppScheme, h.cfg.AppDomain)); err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, nil)
}

func (h *authHandler) ResetPasswordHandler(c *gin.Context) {
	var req schema.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, utils.NewAppError(http.StatusBadRequest, utils.ValidationError, "Invalid JSON payload", err))
		return
	}
	if err := h.validate.Struct(req); err != nil {
		utils.BindErrorResponse(c, utils.NewAppError(http.StatusBadRequest, utils.ValidationError, err.Error(), err))
		return
	}
	if err := h.authService.ResetPassword(c.Request.Context(), req.Token, req.Password); err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, nil)
}

func (h *authHandler) GoogleLoginHandler(c *gin.Context) {
	url := h.authService.GoogleLogin(c.Request.Context())
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *authHandler) GoogleCallbackHandler(c *gin.Context) {
	code := c.Query("code")
	access, refresh, err := h.authService.GoogleCallback(c.Request.Context(), code)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	c.SetCookie("access_token", access, int(h.cfg.AuthAccessTokenExpiration), "/", h.cfg.AppDomain, true, true)
	c.SetCookie("refresh_token", refresh, int(h.cfg.AuthRefreshTokenExpiration), "/", h.cfg.AppDomain, true, true)
	c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s://%s/uptime", h.cfg.AppScheme, h.cfg.AppDomain))
}
