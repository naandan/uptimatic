package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/iancoleman/strcase"
)

func SuccessResponse(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{
		"request_id": getRequestID(c),
		"data":       data,
	})
}

func PaginatedResponse(c *gin.Context, data any, count, limit, page, totalPage int) {
	c.JSON(http.StatusOK, gin.H{
		"request_id": getRequestID(c),
		"data":       data,
		"meta": gin.H{
			"total":       count,
			"limit":       limit,
			"page":        page,
			"total_pages": totalPage,
		},
	})
}

func ErrorResponse(c *gin.Context, status int, code, message string) {
	c.JSON(status, gin.H{
		"request_id": getRequestID(c),
		"error": gin.H{
			"code":    code,
			"message": message,
		},
	})
}

func BindErrorResponse(c *gin.Context, err error) {
	requestID := getRequestID(c)

	if verrs, ok := err.(validator.ValidationErrors); ok {
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

		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"request_id": requestID,
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
