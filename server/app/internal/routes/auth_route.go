package routes

import (
	"uptimatic/internal/handlers"
	"uptimatic/internal/middlewares"
	"uptimatic/internal/utils"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.RouterGroup, h handlers.AuthHandler, jwtUtil *utils.JWTUtil) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", h.RegisterHandler)
		auth.POST("/login", h.LoginHandler)
		auth.POST("/logout", h.LogoutHandler)
		auth.POST("/refresh", h.RefreshHandler)
		auth.GET("/verify", h.VerifyHandler)
		auth.POST("/resend-verification", middlewares.AuthMiddleware(jwtUtil), h.ResendVerificationHandler)
		auth.GET("/resend-verification-ttl", middlewares.AuthMiddleware(jwtUtil), h.ResendVerificationEmailTTLHandler)
		auth.POST("/forgot-password", h.SendPasswordResetEmailHandler)
		auth.POST("/reset-password", h.ResetPasswordHandler)
		auth.GET("/google/login", h.GoogleLoginHandler)
		auth.GET("/google/callback", h.GoogleCallbackHandler)
	}
}
