package routes

import (
	"uptimatic/internal/handlers"
	"uptimatic/internal/middlewares"
	"uptimatic/internal/utils"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.RouterGroup, h handlers.UserHandler, jwtUtil *utils.JWTUtil) {
	users := r.Group("/users")
	users.Use(middlewares.AuthMiddleware(jwtUtil))
	users.GET("/me", h.GetUserHandler)
	users.Use(middlewares.VerifiedMiddleware())
	{
		users.PUT("", h.UpdateUserHandler)
		users.PUT("/change-password", h.ChangePasswordHandler)
		users.POST("/upload-url", h.GetPresignedUrlHandler)
		users.PUT("/update-foto", h.UpdateFotoHandler)
	}
}
