package server

import (
	"context"
	"fmt"
	"time"
	"uptimatic/internal/adapters/google"
	"uptimatic/internal/adapters/minio"
	"uptimatic/internal/auth"
	"uptimatic/internal/config"
	"uptimatic/internal/db"
	"uptimatic/internal/middleware"
	"uptimatic/internal/url"
	"uptimatic/internal/user"
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

	userRepo := user.NewUserRepository()
	urlRepo := url.NewUrlRepository()
	logRepo := url.NewLogRepository()

	authService := auth.NewAuthService(pgsql, userRepo, redis, jwtUtil, asyncClient, googleClient)
	urlService := url.NewUrlService(pgsql, urlRepo, logRepo)
	userService := user.NewUserService(pgsql, userRepo, minio, redis, jwtUtil, asyncClient)

	authHandler := auth.NewAuthHandler(authService, validate, &cfg)
	urlHandler := url.NewURLHandler(urlService, validate)
	userHandler := user.NewUserHandler(userService, validate, &cfg)

	if cfg.AppDebug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.RequestID())

	api := r.Group("/api/v1")
	{
		auth.AuthRoutes(api, authHandler, &jwtUtil)
		user.UserRoutes(api, userHandler, &jwtUtil)
		url.UrlRoutes(api, urlHandler, &jwtUtil)
	}

	addr := ":" + fmt.Sprint(cfg.AppPort)
	if err := r.Run(addr); err != nil {
		utils.Fatal(context.Background(), "Failed to start server", map[string]any{"error": err})
	}
}
