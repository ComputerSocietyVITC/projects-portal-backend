package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ComputerSocietyVITC/projects-portal-backend/internal/config"
	"github.com/ComputerSocietyVITC/projects-portal-backend/internal/handlers"
	"github.com/ComputerSocietyVITC/projects-portal-backend/internal/logger"
	custommiddleware "github.com/ComputerSocietyVITC/projects-portal-backend/internal/middleware"
	"github.com/ComputerSocietyVITC/projects-portal-backend/internal/models"
	"github.com/ComputerSocietyVITC/projects-portal-backend/internal/repository"
	"github.com/ComputerSocietyVITC/projects-portal-backend/internal/service"
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
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPatch, http.MethodDelete},
	}))

	e.GET("/", func(c *echo.Context) error {
		return c.String(http.StatusOK, "Server is up and running!")
	})

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	userRoleRepo := repository.NewUserRoleRepository(db)
	projectMemberRepo := repository.NewProjectMemberRepository(db)
	projectsRepo := repository.NewProjectRepository(db)
	inviteRepo := repository.NewInviteRepository(db)

	// Initialize services
	inviteService := service.NewInviteService(inviteRepo)
	authService := service.NewAuthService(userRepo, userRoleRepo, roleRepo, inviteService)
	userService := service.NewUserService(userRepo)
	projectService := service.NewProjectService(projectsRepo)
	projectMemberService := service.NewProjectMemberService(projectMemberRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService, logger)
	userHandler := handlers.NewUserHandler(userService, logger)
	projectMemberHandler := handlers.NewProjectMemberHandler(projectMemberService, logger)
	projectsHandler := handlers.NewProjectHandler(projectService, logger)
	inviteHandler := handlers.NewInviteHandler(inviteService, logger)

	// Public routes (no authentication required)
	e.POST("/auth/register", authHandler.Register)
	e.POST("/auth/login", authHandler.Login)

	// Protected routes (authentication required)
	api := e.Group("", custommiddleware.AuthMiddleware())

	// User routes - GET all users is group_head only; GET by ID enforces own-profile for members
	api.GET("/users", userHandler.GetAllUsers, custommiddleware.RequireRole("group_head"))
	api.GET("/users/:id", userHandler.GetUserByID)

	// User routes - DELETE only for group heads
	api.DELETE("/users/:id", userHandler.DeleteUser, custommiddleware.RequireRole("group_head"))

	// Project routes - GET all projects is open to all authenticated users; other operations are group_head only
	api.GET("/projects", projectsHandler.GetAllProjects)
	api.GET("/projects/:id", projectsHandler.GetProjectByID)
	api.POST("/projects", projectsHandler.CreateProject, custommiddleware.RequireRole("group_head"))
	api.PUT("/projects/:id", projectsHandler.UpdateProject, custommiddleware.RequireRole("group_head"))
	api.DELETE("/projects/:id", projectsHandler.DeleteProject, custommiddleware.RequireRole("group_head"))

	// Project member routes - PATCH (add user to project) only for group heads
	api.POST("/projects/:project_id/members", projectMemberHandler.AddMember, custommiddleware.RequireRole("group_head"))
	api.DELETE("/projects/:project_id/members/:user_id", projectMemberHandler.RemoveMember, custommiddleware.RequireRole("group_head"))
	api.GET("/projects/:project_id/members", projectMemberHandler.GetMembers)

	// Invite routes - only group heads can create invites
	api.POST("/invites", inviteHandler.CreateInvite, custommiddleware.RequireRole("group_head"))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
		logger.Fatal("server error: %v", zap.Error(err))
	}
}
