package server

import (
	"context"
	"fmt"
	"time"
	"uptimatic/internal/config"
	"uptimatic/internal/db"
	"uptimatic/internal/google"
	"uptimatic/internal/handlers"
	"uptimatic/internal/middlewares"
	"uptimatic/internal/minio"
	"uptimatic/internal/repositories"
	"uptimatic/internal/routes"
	"uptimatic/internal/services"
	"uptimatic/internal/utils"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func Start() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	utils.InitLogger(cfg.AppLogLevel)
	_ = utils.InitSentry(cfg.SentryDSN)
	defer sentry.Flush(2 * time.Second)

	pgsql := db.NewPostgresClient(&cfg)
	redis := db.NewRedisClient(&cfg)
	asyncClient := db.NewAsynqClient(&cfg)

	jwtUtil := utils.NewJWTUtil(cfg.AuthJWTSecret, cfg.AuthAccessTokenExpiration, cfg.AuthRefreshTokenExpiration)
	googleClient := google.NewGoogleClient(&cfg)
	validate := validator.New()

	minio, err := minio.NewMinioUtil(context.Background(), &cfg)
	if err != nil {
		utils.Fatal(context.Background(), "Failed to connect minio", map[string]any{"error": err})
	}

	userRepo := repositories.NewUserRepository()
	urlRepo := repositories.NewUrlRepository()
	logRepo := repositories.NewLogRepository()

	authService := services.NewAuthService(pgsql, userRepo, redis, jwtUtil, asyncClient, googleClient)
	urlService := services.NewUrlService(pgsql, urlRepo, logRepo)
	userService := services.NewUserService(pgsql, userRepo, minio, redis, jwtUtil, asyncClient)

	authHandler := handlers.NewAuthHandler(authService, validate, &cfg)
	urlHandler := handlers.NewURLHandler(urlService, validate)
	userHandler := handlers.NewUserHandler(userService, validate, &cfg)

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
		routes.AuthRoutes(api, authHandler, &jwtUtil)
		routes.UserRoutes(api, userHandler, &jwtUtil)
		routes.UrlRoutes(api, urlHandler, &jwtUtil)
	}

	addr := ":" + fmt.Sprint(cfg.AppPort)
	if err := r.Run(addr); err != nil {
		utils.Fatal(context.Background(), "Failed to start server", map[string]any{"error": err})
	}
}
