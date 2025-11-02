package url

import (
	"net/http"
	"strconv"
	"uptimatic/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
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
	urlService URLService
	validate   *validator.Validate
}

func NewURLHandler(urlService URLService, validate *validator.Validate) URLHandler {
	return &urlHandler{urlService, validate}
}

func (h *urlHandler) CreateHandler(c *gin.Context) {
	var urlRequest UrlRequest
	if err := c.ShouldBindJSON(&urlRequest); err != nil {
		utils.ErrorResponse(c, utils.NewAppError(http.StatusBadRequest, utils.ValidationError, "Invalid JSON payload", err))
		return
	}
	if err := h.validate.Struct(urlRequest); err != nil {
		utils.BindErrorResponse(c, utils.NewAppError(http.StatusBadRequest, utils.ValidationError, err.Error(), err))
		return
	}

	userId := c.GetUint("user_id")
	urlResponse, err := h.urlService.Create(c.Request.Context(), &urlRequest, userId)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, urlResponse)
}

func (h *urlHandler) UpdateHandler(c *gin.Context) {
	id := c.Param("id")
	publicID, err := uuid.Parse(id)
	if err != nil {
		utils.ErrorResponse(c, utils.NewAppError(http.StatusBadRequest, utils.ValidationError, err.Error(), err))
		return
	}

	var urlRequest UrlRequest
	if err := c.ShouldBindJSON(&urlRequest); err != nil {
		utils.ErrorResponse(c, utils.NewAppError(http.StatusBadRequest, utils.ValidationError, "Invalid JSON payload", err))
		return
	}
	if err := h.validate.Struct(urlRequest); err != nil {
		utils.BindErrorResponse(c, utils.NewAppError(http.StatusBadRequest, utils.ValidationError, err.Error(), err))
		return
	}

	urlResponse, errSvc := h.urlService.Update(c.Request.Context(), &urlRequest, publicID)
	if errSvc != nil {
		utils.ErrorResponse(c, errSvc)
		return
	}

	utils.SuccessResponse(c, urlResponse)
}

func (h *urlHandler) DeleteHandler(c *gin.Context) {
	id := c.Param("id")
	publicID, err := uuid.Parse(id)
	if err != nil {
		utils.ErrorResponse(c, utils.NewAppError(http.StatusBadRequest, utils.ValidationError, err.Error(), err))
		return
	}

	if err := h.urlService.Delete(c.Request.Context(), publicID); err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, nil)
}

func (h *urlHandler) GetHandler(c *gin.Context) {
	id := c.Param("id")
	publicID, err := uuid.Parse(id)
	if err != nil {
		utils.ErrorResponse(c, utils.NewAppError(http.StatusBadRequest, utils.ValidationError, err.Error(), err))
		return
	}

	urlResponse, errSvc := h.urlService.FindByID(c.Request.Context(), publicID)
	if errSvc != nil {
		utils.ErrorResponse(c, errSvc)
		return
	}
	utils.SuccessResponse(c, urlResponse)
}

func (h *urlHandler) ListHandler(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		utils.ErrorResponse(c, utils.NewAppError(http.StatusBadRequest, utils.ValidationError, "Invalid page", err))
		return
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		utils.ErrorResponse(c, utils.NewAppError(http.StatusBadRequest, utils.ValidationError, "Invalid limit", err))
		return
	}

	var active *bool
	if a := c.Query("active"); a != "" {
		switch a {
		case "active":
			active = new(bool)
			*active = true
		case "inactive":
			active = new(bool)
			*active = false
		default:
			active = nil
		}
	}

	searchLabel := c.Query("q")
	sortBy := c.DefaultQuery("sort", "label")

	urls, count, errSvc := h.urlService.ListByUserID(c.Request.Context(), c.GetUint("user_id"), page, limit, active, searchLabel, sortBy)
	if errSvc != nil {
		utils.ErrorResponse(c, errSvc)
		return
	}

	utils.PaginatedResponse(c, urls, count, limit, page, (count+limit-1)/limit)
}

func (h *urlHandler) GetUptimeStats(c *gin.Context) {
	mode := c.DefaultQuery("mode", "day")
	dateStr := c.DefaultQuery("date", "")
	id := c.Param("id")

	if mode != "day" && mode != "month" && mode != "year" {
		utils.ErrorResponse(c, utils.NewAppError(http.StatusBadRequest, utils.ValidationError, "Invalid mode", nil))
		return
	}

	idUUID, err := uuid.Parse(id)
	if err != nil {
		utils.ErrorResponse(c, utils.NewAppError(http.StatusBadRequest, utils.ValidationError, err.Error(), err))
		return
	}

	stats, errSvc := h.urlService.GetUptimeStats(c.Request.Context(), idUUID, mode, dateStr)
	if errSvc != nil {
		utils.ErrorResponse(c, errSvc)
		return
	}

	utils.SuccessResponse(c, stats)
}
