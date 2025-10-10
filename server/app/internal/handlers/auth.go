package handlers

import (
	"net/http"
	"uptimatic/internal/config"
	"uptimatic/internal/schema"
	service "uptimatic/internal/services"
	"uptimatic/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthHandler interface {
	RegisterHandler(c *gin.Context)
	LoginHandler(c *gin.Context)
	LogoutHandler(c *gin.Context)
	RefreshHandler(c *gin.Context)
}

type authHandler struct {
	authService service.AuthService
	validate    *validator.Validate
	cfg         *config.Config
}

func NewAuthHandler(authService service.AuthService, validate *validator.Validate, cfg *config.Config) AuthHandler {
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

	user, err := h.authService.Register(req.Email, req.Password, "")
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

	access, refresh, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
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

func (h *authHandler) LogoutHandler(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, utils.Unauthorized, "Missing refresh token")
		return
	}

	if err := h.authService.Logout(refreshToken); err != nil {
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

	access, refresh, err := h.authService.Refresh(refreshToken)
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
