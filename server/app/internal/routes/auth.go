package routes

import (
	"uptimatic/internal/handlers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.RouterGroup, h handlers.AuthHandler) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", h.RegisterHandler)
		auth.POST("/login", h.LoginHandler)
		auth.POST("/logout", h.LogoutHandler)
		auth.POST("/refresh", h.RefreshHandler)
	}
}
