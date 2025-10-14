package routes

import (
	"uptimatic/internal/handlers"
	"uptimatic/internal/middlewares"
	"uptimatic/internal/utils"

	"github.com/gin-gonic/gin"
)

func UrlRoutes(r *gin.RouterGroup, h handlers.URLHandler, jwtUtil *utils.JWTUtil) {
	urls := r.Group("/urls")
	urls.Use(middlewares.AuthMiddleware(jwtUtil))
	urls.Use(middlewares.VerifiedMiddleware())
	{
		urls.POST("", h.CreateHandler)
		urls.GET("", h.ListHandler)
		urls.GET("/:id", h.GetHandler)
		urls.PUT("/:id", h.UpdateHandler)
		urls.DELETE("/:id", h.DeleteHandler)
		urls.GET("/:id/stats", h.GetUptimeStats)
	}
}
