package server

import (
	"fmt"
	"uptimatic/internal/config"
	"uptimatic/internal/db"
	"uptimatic/internal/handlers"
	"uptimatic/internal/middlewares"
	"uptimatic/internal/repositories"
	"uptimatic/internal/routes"
	"uptimatic/internal/services"
	"uptimatic/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func Start() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	utils.InitLogger(cfg.AppLogLevel)

	pgsql := db.NewPostgresClient(&cfg)
	redis := db.NewRedisClient(&cfg)
	asyncClient := db.NewAsynqClient(&cfg)

	jwtUtil := utils.NewJWTUtil(cfg.AuthJWTSecret, cfg.AuthAccessTokenExpiration, cfg.AuthRefreshTokenExpiration)
	validate := validator.New()

	userRepo := repositories.NewUserRepository()

	authService := services.NewAuthService(pgsql, userRepo, redis, jwtUtil, asyncClient)

	authHandler := handlers.NewAuthHandler(authService, validate, &cfg)

	if cfg.AppDebug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middlewares.RequestID())

	api := r.Group("/api/v1")
	{
		routes.AuthRoutes(api, authHandler)
	}

	addr := ":" + fmt.Sprint(cfg.AppPort)
	if err := r.Run(addr); err != nil {
		utils.Error(nil, "failed to start server", map[string]any{"error": err})
	}
}
