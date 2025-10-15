package middlewares

import (
	"net/http"
	"uptimatic/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func VerifiedMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		claimsVal, exists := c.Get("claims")
		if !exists {
			utils.ErrorResponse(c, http.StatusUnauthorized, utils.InvalidToken, "Missing token claims")
			c.Abort()
			return
		}

		var verified bool
		switch v := claimsVal.(type) {
		case jwt.MapClaims:
			if val, ok := v["verified"].(bool); ok {
				verified = val
			}
		case map[string]interface{}:
			if val, ok := v["verified"].(bool); ok {
				verified = val
			}
		default:
			utils.ErrorResponse(c, http.StatusUnauthorized, utils.InvalidToken, "Invalid claims type")
			c.Abort()
			return
		}

		if !verified {
			utils.ErrorResponse(c, http.StatusForbidden, utils.ForbiddenAction, "Email not verified")
			c.Abort()
			return
		}

		c.Next()
	}
}
