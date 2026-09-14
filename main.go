package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/ComputerSocietyVITC/projects-portal-backend/internal/config"
	"github.com/ComputerSocietyVITC/projects-portal-backend/internal/handlers"
	"github.com/ComputerSocietyVITC/projects-portal-backend/internal/logger"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"go.uber.org/zap"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("failed to load .env file: %v", err)
	}

	logger, err := logger.Init()
	if err != nil {
		log.Fatalf("failed to initialize zap logger: %v", err)
	}
	zap.ReplaceGlobals(logger)

	db, err := config.ConnectDatabase()
	if err != nil {
		logger.Fatal("failed to connect to database: %v", zap.Error(err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Fatal("failed to get sql db: %v", zap.Error(err))
	}
	if err := sqlDB.Ping(); err != nil {
		logger.Fatal("failed to ping database: %v", zap.Error(err))
	}

	redisDB := config.ConnectRedis()

	_, err = redisDB.Ping(context.Background()).Result()
	if err != nil {
		logger.Fatal("failed to connect to redis store: %v", zap.Error(err))
	}

	e := echo.New()
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c *echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info("request",
				zap.String("URI", v.URI),
				zap.Int("status", v.Status),
			)

			return nil
		},
	}))

	e.GET("/", func(c *echo.Context) error {
		return c.String(http.StatusOK, "Server is up and running!")
	})

	handlers.InitRoutes(e, db, redisDB, logger)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
		logger.Fatal("server error: %v", zap.Error(err))
	}
}
