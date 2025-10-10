package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

const (
	Required        = "REQUIRED"
	InvalidType     = "INVALID_TYPE"
	InvalidFormat   = "INVALID_FORMAT"
	MinLength       = "MIN_LENGTH"
	MaxLength       = "MAX_LENGTH"
	ValidationError = "VALIDATION_ERROR"

	NotFound = "NOT_FOUND"

	Unauthorized    = "UNAUTHORIZED"
	InvalidToken    = "INVALID_TOKEN"
	ForbiddenAction = "FORBIDDEN_ACTION"

	InternalError = "INTERNAL_ERROR"
	Conflict      = "CONFLICT"
)

func SuccessResponse(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{
		"request_id": getRequestID(c),
		"status":     "success",
		"data":       data,
	})
}

func ErrorResponse(c *gin.Context, status int, code, message string) {
	c.JSON(status, gin.H{
		"request_id": getRequestID(c),
		"status":     "error",
		"error": gin.H{
			"code":    code,
			"message": message,
		},
	})
}

func BindErrorResponse(c *gin.Context, err error) {
	requestID := getRequestID(c)

	if verrs, ok := err.(validator.ValidationErrors); ok {
		errors := make(map[string][]string)

		for _, e := range verrs {
			field := e.Field()
			switch e.Tag() {
			case "required":
				errors[field] = append(errors[field], Required)
			case "email":
				errors[field] = append(errors[field], InvalidFormat)
			case "min":
				errors[field] = append(errors[field], MinLength)
			case "max":
				errors[field] = append(errors[field], MaxLength)
			default:
				errors[field] = append(errors[field], InvalidType)
			}
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"request_id": requestID,
			"status":     "error",
			"error": gin.H{
				"code":    ValidationError,
				"message": "Payload validation failed",
				"fields":  errors,
			},
		})
		return
	}

	ErrorResponse(c, http.StatusBadRequest, ValidationError, err.Error())
}

func getRequestID(c *gin.Context) string {
	if reqID, exists := c.Get("request_id"); exists {
		if idStr, ok := reqID.(string); ok {
			return idStr
		}
	}
	return uuid.New().String()
}
