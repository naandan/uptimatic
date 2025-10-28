package middleware

import (
	"context"
	"uptimatic/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.Request.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		ctx := context.WithValue(c.Request.Context(), utils.TraceKey, requestID)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
