package url

import (
	"uptimatic/internal/middleware"
	"uptimatic/internal/utils"

	"github.com/gin-gonic/gin"
)

func UrlRoutes(r *gin.RouterGroup, h URLHandler, jwtUtil *utils.JWTUtil) {
	urls := r.Group("/urls")
	urls.Use(middleware.AuthMiddleware(jwtUtil))
	urls.Use(middleware.VerifiedMiddleware())
	{
		urls.POST("", h.CreateHandler)
		urls.GET("", h.ListHandler)
		urls.GET("/:id", h.GetHandler)
		urls.PUT("/:id", h.UpdateHandler)
		urls.DELETE("/:id", h.DeleteHandler)
		urls.GET("/:id/stats", h.GetUptimeStats)
	}
}
