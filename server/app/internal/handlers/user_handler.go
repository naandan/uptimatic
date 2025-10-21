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

type UserHandler interface {
	UpdateUserHandler(c *gin.Context)
	GetUserHandler(c *gin.Context)
	ChangePasswordHandler(c *gin.Context)
	GetPresignedUrlHandler(c *gin.Context)
	UpdateFotoHandler(c *gin.Context)
}

type userHandler struct {
	userService services.UserService
	validate    *validator.Validate
	cfg         *config.Config
}

func NewUserHandler(userService services.UserService, validate *validator.Validate, cfg *config.Config) UserHandler {
	return &userHandler{userService, validate, cfg}
}

func (h *userHandler) UpdateUserHandler(c *gin.Context) {
	var req schema.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, utils.NewAppError(http.StatusBadRequest, utils.ValidationError, "Invalid JSON payload", err))
		return
	}
	if err := h.validate.Struct(req); err != nil {
		utils.BindErrorResponse(c, utils.NewAppError(http.StatusBadRequest, utils.ValidationError, err.Error(), err))
		return
	}
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.SetCookie("access_token", "", -1, "/", h.cfg.AppDomain, true, true)
		c.SetCookie("refresh_token", "", -1, "/", h.cfg.AppDomain, true, true)
		utils.ErrorResponse(c, utils.NewAppError(http.StatusUnauthorized, utils.InvalidCredentials, "Missing refresh token", nil))
		return
	}

	user, changed, errSvc := h.userService.Update(c.Request.Context(), c.GetUint("user_id"), req.Name, req.Email, fmt.Sprintf("%s://%s", h.cfg.AppScheme, h.cfg.AppDomain), refreshToken)
	if errSvc != nil {
		utils.ErrorResponse(c, errSvc)
		return
	}

	if changed != nil && changed["is_email_changed"].(bool) {
		c.SetCookie("access_token", changed["access_token"].(string), int(h.cfg.AuthAccessTokenExpiration.Seconds()), "/", h.cfg.AppDomain, true, true)
		c.SetCookie("refresh_token", changed["refresh_token"].(string), int(h.cfg.AuthRefreshTokenExpiration.Seconds()), "/", h.cfg.AppDomain, true, true)
	}

	utils.SuccessResponse(c, user)
}

func (h *userHandler) GetUserHandler(c *gin.Context) {
	user, err := h.userService.GetUser(c.Request.Context(), c.GetUint("user_id"))
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, user)
}

func (h *userHandler) ChangePasswordHandler(c *gin.Context) {
	var req schema.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, utils.NewAppError(http.StatusBadRequest, utils.ValidationError, "Invalid JSON payload", err))
		return
	}
	if err := h.validate.Struct(req); err != nil {
		utils.BindErrorResponse(c, utils.NewAppError(http.StatusBadRequest, utils.ValidationError, err.Error(), err))
		return
	}
	err := h.userService.ChangePassword(c.Request.Context(), c.GetUint("user_id"), req.OldPassword, req.NewPassword)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, nil)
}

func (h *userHandler) GetPresignedUrlHandler(c *gin.Context) {
	var req schema.PresignedUrlRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, utils.NewAppError(http.StatusBadRequest, utils.ValidationError, "Invalid JSON payload", err))
		return
	}

	if err := h.validate.Struct(req); err != nil {
		utils.BindErrorResponse(c, utils.NewAppError(http.StatusBadRequest, utils.ValidationError, err.Error(), err))
		return
	}

	url, fileName, err := h.userService.GetPresignedUrl(c.Request.Context(), req.FileName, req.ContentType)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, gin.H{
		"file_name":     fileName,
		"presigned_url": url,
	})
}

func (h *userHandler) UpdateFotoHandler(c *gin.Context) {
	var req schema.UpdateFotoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, utils.NewAppError(http.StatusBadRequest, utils.ValidationError, "Invalid JSON payload", err))
		return
	}

	if err := h.validate.Struct(req); err != nil {
		utils.BindErrorResponse(c, utils.NewAppError(http.StatusBadRequest, utils.ValidationError, err.Error(), err))
		return
	}

	url, err := h.userService.UpdateFoto(c.Request.Context(), c.GetUint("user_id"), req.FileName)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, gin.H{
		"url": url,
	})
}
