package handlers

import (
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
		utils.ErrorResponse(c, http.StatusBadRequest, utils.ValidationError, "Invalid JSON payload")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		utils.BindErrorResponse(c, err)
		return
	}

	user, err := h.authService.Register(c.Request.Context(), req.Email, req.Password, h.cfg.AppDomain)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.InternalError, err.Error())
		return
	}

	utils.SuccessResponse(c, user)
}

func (h *authHandler) LoginHandler(c *gin.Context) {
	var req schema.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, utils.ValidationError, "Invalid JSON payload")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		utils.BindErrorResponse(c, err)
		return
	}

	access, refresh, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, utils.ValidationError, err.Error())
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
		utils.ErrorResponse(c, http.StatusUnauthorized, utils.Unauthorized, "Missing refresh token")
		return
	}

	if err := h.authService.Logout(c.Request.Context(), refreshToken); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.InternalError, err.Error())
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
		utils.ErrorResponse(c, http.StatusUnauthorized, utils.Unauthorized, "Missing refresh token")
		return
	}

	access, refresh, err := h.authService.Refresh(c.Request.Context(), refreshToken)
	if err != nil {
		c.SetCookie("access_token", "", -1, "/", h.cfg.AppDomain, true, true)
		c.SetCookie("refresh_token", "", -1, "/", h.cfg.AppDomain, true, true)
		utils.ErrorResponse(c, http.StatusUnauthorized, utils.Unauthorized, err.Error())
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
		utils.ErrorResponse(c, http.StatusUnauthorized, utils.Unauthorized, "User not authenticated")
		return
	}

	userIdUint, ok := userId.(uint)
	if !ok {
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.InternalError, "Invalid user ID type")
		return
	}

	user, err := h.authService.Profile(c.Request.Context(), userIdUint)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.InternalError, err.Error())
		return
	}

	utils.SuccessResponse(c, user)
}

func (h *authHandler) VerifyHandler(c *gin.Context) {
	token := c.Query("token")
	if err := h.authService.VerifyEmail(c.Request.Context(), token); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.InternalError, err.Error())
		return
	}
	utils.SuccessResponse(c, nil)
}

func (h *authHandler) ResendVerificationHandler(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		utils.ErrorResponse(c, http.StatusUnauthorized, utils.Unauthorized, "User not authenticated")
		return
	}

	ttl, err := h.authService.ResendVerificationEmail(c.Request.Context(), userID, h.cfg.AppDomain)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.InternalError, err.Error())
		return
	}
	utils.SuccessResponse(c, gin.H{"ttl": ttl})
}

func (h *authHandler) ResendVerificationEmailTTLHandler(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		utils.ErrorResponse(c, http.StatusUnauthorized, utils.Unauthorized, "User not authenticated")
		return
	}

	ttl, err := h.authService.ResendVerificationEmailTTL(c.Request.Context(), userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.InternalError, err.Error())
		return
	}
	utils.SuccessResponse(c, gin.H{"ttl": ttl})
}

func (h *authHandler) SendPasswordResetEmailHandler(c *gin.Context) {
	var req schema.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, utils.ValidationError, err.Error())
		return
	}
	if err := h.validate.Struct(req); err != nil {
		utils.BindErrorResponse(c, err)
		return
	}
	if err := h.authService.SendPasswordResetEmail(c.Request.Context(), req.Email, h.cfg.AppDomain); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.InternalError, err.Error())
		return
	}
	utils.SuccessResponse(c, nil)
}

func (h *authHandler) ResetPasswordHandler(c *gin.Context) {
	var req schema.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, utils.ValidationError, err.Error())
		return
	}
	if err := h.validate.Struct(req); err != nil {
		utils.BindErrorResponse(c, err)
		return
	}
	if err := h.authService.ResetPassword(c.Request.Context(), req.Token, req.Password); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.InternalError, err.Error())
		return
	}
	utils.SuccessResponse(c, nil)
}
