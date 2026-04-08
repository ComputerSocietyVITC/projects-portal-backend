package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ComputerSocietyVITC/projects-portal-backend/internal/config"
	"github.com/ComputerSocietyVITC/projects-portal-backend/internal/handlers"
	"github.com/ComputerSocietyVITC/projects-portal-backend/internal/logger"
	"github.com/ComputerSocietyVITC/projects-portal-backend/internal/models"
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

	if err := db.AutoMigrate(&models.Project{}); err != nil {
		logger.Fatal("failed to auto-migrate database: %v", zap.Error(err))
	}

	if err := sqlDB.Ping(); err != nil {
		logger.Fatal("failed to ping database: %v", zap.Error(err))
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

	e.GET("/users", func(c *echo.Context) error {
		handlers := &handlers.UserHandler{DB: &config.Database{DB: db}}
		return handlers.GetUsers(c)
	})

	e.GET("/users/:id", func(c *echo.Context) error {
		handlers := &handlers.UserHandler{DB: &config.Database{DB: db}}
		return handlers.GetUserByID(c)
	})

	e.POST("/users", func(c *echo.Context) error {
		handlers := &handlers.UserHandler{DB: &config.Database{DB: db}}
		return handlers.CreateUser(c)
	})

	e.GET("/projects", func(c *echo.Context) error {
		handlers := &handlers.ProjectHandler{DB: &config.Database{DB: db}}
		return handlers.GetProjects(c)
	})
	e.GET("/projects/:id", func(c *echo.Context) error {
		handlers := &handlers.ProjectHandler{DB: &config.Database{DB: db}}
		return handlers.GetProjectByID(c)
	})

	e.POST("/projects", func(c *echo.Context) error {
		handlers := &handlers.ProjectHandler{DB: &config.Database{DB: db}}
		return handlers.CreateProject(c)
	})

	e.PATCH("/projects/:id", func(c *echo.Context) error {
		handlers := &handlers.ProjectHandler{DB: &config.Database{DB: db}}
		return handlers.UpdateProject(c)
	})

	e.DELETE("/projects/:id", func(c *echo.Context) error {
		handlers := &handlers.ProjectHandler{DB: &config.Database{DB: db}}
		return handlers.DeleteProject(c)
	})

	e.POST("/projects/:id/members", func(c *echo.Context) error {
		handlers := &handlers.ProjectMemberHandler{DB: &config.Database{DB: db}}
		return handlers.AddMember(c)
	})

	e.DELETE("/projects/:id/members/:member_id", func(c *echo.Context) error {
		handlers := &handlers.ProjectMemberHandler{DB: &config.Database{DB: db}}
		return handlers.RemoveMember(c)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
		logger.Fatal("server error: %v", zap.Error(err))
	}
}
