package user

import (
	"uptimatic/internal/middleware"
	"uptimatic/internal/utils"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.RouterGroup, h UserHandler, jwtUtil *utils.JWTUtil) {
	users := r.Group("/users")
	users.Use(middleware.AuthMiddleware(jwtUtil))
	users.GET("/me", h.GetUserHandler)
	users.Use(middleware.VerifiedMiddleware())
	{
		users.PUT("", h.UpdateUserHandler)
		users.PUT("/change-password", h.ChangePasswordHandler)
		users.POST("/upload-url", h.GetPresignedUrlHandler)
		users.PUT("/update-foto", h.UpdateFotoHandler)
	}
}
