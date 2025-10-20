package utils

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/iancoleman/strcase"
)

const (
	ValidationError = "VALIDATION_ERROR"

	Unauthorized       = "UNAUTHORIZED"
	InvalidCredentials = "INVALID_CREDENTIALS"
	InvalidToken       = "INVALID_TOKEN"
	ForbiddenAction    = "FORBIDDEN_ACTION"
	AccountLocked      = "ACCOUNT_LOCKED"
	TooManyRequests    = "TOO_MANY_REQUESTS"

	NotFound = "NOT_FOUND"
	Conflict = "CONFLICT"

	InternalError      = "INTERNAL_ERROR"
	ServiceUnavailable = "SERVICE_UNAVAILABLE"
	Timeout            = "TIMEOUT"
)

const (
	Required      = "REQUIRED"
	InvalidType   = "INVALID_TYPE"
	InvalidFormat = "INVALID_FORMAT"
	MinLength     = "MIN_LENGTH"
	MaxLength     = "MAX_LENGTH"
	EnumValue     = "ENUM_VALUE"
	Mismatch      = "MISMATCH"
	Unique        = "UNIQUE"
	MinValue      = "MIN_VALUE"
	MaxValue      = "MAX_VALUE"
)

func SuccessResponse(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{
		"request_id": getRequestID(c.Request.Context()),
		"data":       data,
	})
}

func PaginatedResponse(c *gin.Context, data any, count, limit, page, totalPage int) {
	c.JSON(http.StatusOK, gin.H{
		"request_id": getRequestID(c.Request.Context()),
		"data":       data,
		"meta": gin.H{
			"total":       count,
			"limit":       limit,
			"page":        page,
			"total_pages": totalPage,
		},
	})
}

func ErrorResponse(c *gin.Context, appErr *AppError) {
	resp := gin.H{
		"request_id": getRequestID(c.Request.Context()),
		"error": gin.H{
			"code":    appErr.Code,
			"message": appErr.Message,
		},
	}
	if appErr.Fields != nil {
		resp["error"].(gin.H)["fields"] = appErr.Fields
	}
	c.JSON(appErr.Status, resp)
}

func BindErrorResponse(c *gin.Context, appErr *AppError) {
	requestID := getRequestID(c.Request.Context())

	if verrs, ok := appErr.Err.(validator.ValidationErrors); ok {
		errors := make(map[string][]map[string]interface{})

		for _, e := range verrs {
			fieldName := strcase.ToSnake(e.Field())

			switch e.Tag() {
			case "required":
				errors[fieldName] = append(errors[fieldName], map[string]interface{}{
					"code": Required,
				})
			case "email":
				errors[fieldName] = append(errors[fieldName], map[string]interface{}{
					"code": InvalidFormat,
				})
			case "min":
				errors[fieldName] = append(errors[fieldName], map[string]interface{}{
					"code":  MinLength,
					"param": e.Param(),
				})
			case "max":
				errors[fieldName] = append(errors[fieldName], map[string]interface{}{
					"code":  MaxLength,
					"param": e.Param(),
				})
			case "oneof":
				errors[fieldName] = append(errors[fieldName], map[string]interface{}{
					"code":  EnumValue,
					"param": e.Param(),
				})
			default:
				errors[fieldName] = append(errors[fieldName], map[string]interface{}{
					"code": InvalidType,
				})
			}
		}

		c.JSON(appErr.Status, gin.H{
			"request_id": requestID,
			"error": gin.H{
				"code":    appErr.Code,
				"message": "Payload validation failed",
				"fields":  errors,
			},
		})
		return
	}

	ErrorResponse(c, appErr)
}

func getRequestID(ctx context.Context) string {
	if reqID, ok := ctx.Value(TraceKey).(string); ok {
		return reqID
	}
	return uuid.New().String()
}
