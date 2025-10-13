package handlers

import (
	"net/http"
	"strconv"
	"uptimatic/internal/schema"
	"uptimatic/internal/services"
	"uptimatic/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type URLHandler interface {
	CreateHandler(c *gin.Context)
	UpdateHandler(c *gin.Context)
	DeleteHandler(c *gin.Context)
	GetHandler(c *gin.Context)
	ListHandler(c *gin.Context)
	GetUptimeStats(c *gin.Context)
}

type urlHandler struct {
	urlService services.URLService
	validate   *validator.Validate
}

func NewURLHandler(urlService services.URLService, validate *validator.Validate) URLHandler {
	return &urlHandler{urlService, validate}
}

func (h *urlHandler) CreateHandler(c *gin.Context) {
	var urlRequest schema.UrlRequest
	if err := c.ShouldBindJSON(&urlRequest); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, utils.ValidationError, err.Error())
		return
	}
	if err := h.validate.Struct(urlRequest); err != nil {
		utils.BindErrorResponse(c, err)
		return
	}

	userId := c.GetUint("user_id")
	urlResponse, err := h.urlService.Create(&urlRequest, userId)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.InternalError, err.Error())
		return
	}

	utils.SuccessResponse(c, urlResponse)
}

func (h *urlHandler) UpdateHandler(c *gin.Context) {
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, utils.ValidationError, err.Error())
		return
	}

	var urlRequest schema.UrlRequest
	if err := c.ShouldBindJSON(&urlRequest); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, utils.ValidationError, err.Error())
		return
	}
	if err := h.validate.Struct(urlRequest); err != nil {
		utils.BindErrorResponse(c, err)
		return
	}

	urlResponse, err := h.urlService.Update(&urlRequest, uint(idUint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.InternalError, err.Error())
		return
	}
	c.JSON(http.StatusOK, urlResponse)
}

func (h *urlHandler) DeleteHandler(c *gin.Context) {
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, utils.ValidationError, err.Error())
		return
	}

	if err := h.urlService.Delete(uint(idUint)); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.InternalError, err.Error())
		return
	}

	utils.SuccessResponse(c, nil)
}

func (h *urlHandler) GetHandler(c *gin.Context) {
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, utils.ValidationError, err.Error())
		return
	}

	urlResponse, err := h.urlService.FindByID(uint(idUint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.InternalError, err.Error())
		return
	}
	utils.SuccessResponse(c, urlResponse)
}

func (h *urlHandler) ListHandler(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, utils.ValidationError, err.Error())
		return
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, utils.ValidationError, err.Error())
		return
	}
	urls, count, err := h.urlService.ListByUserID(c.GetUint("user_id"), page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.InternalError, err.Error())
		return
	}

	utils.PaginatedResponse(c, urls, count, limit, page, (count+limit-1)/limit)
}

func (h *urlHandler) GetUptimeStats(c *gin.Context) {
	mode := c.DefaultQuery("mode", "day")
	offsetStr := c.DefaultQuery("offset", "0")
	id := c.Param("id")

	if mode != "day" && mode != "month" {
		utils.ErrorResponse(c, http.StatusBadRequest, utils.ValidationError, "invalid mode")
		return
	}

	offset, _ := strconv.Atoi(offsetStr)
	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, utils.ValidationError, err.Error())
		return
	}

	stats, err := h.urlService.GetUptimeStats(uint(idUint), mode, offset)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.InternalError, err.Error())
		return
	}

	utils.SuccessResponse(c, stats)
}
