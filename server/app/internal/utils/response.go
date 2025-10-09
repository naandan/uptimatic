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
	MinValue        = "MIN_VALUE"
	MaxValue        = "MAX_VALUE"
	Unique          = "UNIQUE"
	EnumValue       = "ENUM_VALUE"
	Mismatch        = "MISMATCH"
	ValidationError = "VALIDATION_ERROR"

	NotFound = "NOT_FOUND"

	Unauthorized    = "UNAUTHORIZED"
	InvalidToken    = "INVALID_TOKEN"
	ForbiddenAction = "FORBIDDEN_ACTION"
	AccountLocked   = "ACCOUNT_LOCKED"
	TooManyRequests = "TOO_MANY_REQUESTS"

	InternalError      = "INTERNAL_ERROR"
	ServiceUnavailable = "SERVICE_UNAVAILABLE"
	Timeout            = "TIMEOUT"
	Conflict           = "CONFLICT"
)

func SuccessResponse(c *gin.Context, data any) {
	requestID := getRequestID(c)
	c.JSON(http.StatusOK, gin.H{
		"request_id": requestID,
		"data":       data,
	})
}

func ErrorResponse(c *gin.Context, status int, code, message string) {
	requestID := getRequestID(c)
	c.JSON(status, gin.H{
		"request_id": requestID,
		"error": gin.H{
			"message": message,
			"code":    code,
		},
	})
}

func BindErrorResponse(c *gin.Context, err error) {
	requestID := getRequestID(c)
	errors := make(map[string][]string)

	if verrs, ok := err.(validator.ValidationErrors); ok {
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
			"error": gin.H{
				"message": "Payload Validation Failed",
				"code":    ValidationError,
				"fields":  errors,
			},
		})
		return
	}

	ErrorResponse(c, http.StatusBadRequest, ValidationError, err.Error())
}

func getRequestID(c *gin.Context) string {
	reqID, exists := c.Get("request_id")
	if !exists {
		return uuid.New().String()
	}
	return reqID.(string)
}
