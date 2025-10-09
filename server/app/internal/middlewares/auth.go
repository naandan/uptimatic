package middlewares

import (
	"net/http"
	"strings"
	"uptimatic/internal/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtUtil *utils.JWTUtil) gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string

		if cookie, err := c.Cookie("access_token"); err == nil {
			token = cookie
		}

		if token == "" {
			authHeader := c.GetHeader("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				token = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		if token == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, utils.Unauthorized, "Missing access token")
			c.Abort()
			return
		}

		claims, err := jwtUtil.ValidateToken(token)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, utils.InvalidToken, "Invalid or expired token")
			c.Abort()
			return
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			utils.ErrorResponse(c, http.StatusUnauthorized, utils.InvalidToken, "Invalid token payload")
			c.Abort()
			return
		}

		c.Set("user_id", uint64(userID))
		c.Next()
	}
}
