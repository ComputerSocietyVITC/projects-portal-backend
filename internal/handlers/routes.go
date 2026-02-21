package handlers

import (
	"github.com/ComputerSocietyVITC/projects-portal-backend/internal/middleware"
	"github.com/ComputerSocietyVITC/projects-portal-backend/internal/repository"
	"github.com/ComputerSocietyVITC/projects-portal-backend/internal/service"
	"github.com/labstack/echo/v5"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func InitRoutes(e *echo.Echo, db *gorm.DB, redis *redis.Client, logger *zap.Logger) {
	inviteRepo := &repository.InviteRepo{DB: db, Logger: logger}
	authRepo := &repository.AuthRepo{DB: db, Logger: logger}
	redisRepo := &repository.RedisRepo{Redis: redis, Logger: logger}

	inviteService := &service.InviteService{Repo: inviteRepo}
	authService := &service.AuthService{Repo: authRepo, RedisRepo: redisRepo, InviteService: inviteService}

	handlers := NewHandlers(authService, inviteService, logger)

	authGroup := e.Group("/auth")
	authGroup.POST("/login", handlers.Auth.Login)
	authGroup.POST("/refresh", handlers.Auth.Refresh)
	authGroup.POST("/register", handlers.Auth.Register)
	authGroup.POST("/logout", handlers.Auth.Logout, middleware.LoggedIn)

	inviteGroup := e.Group("/invites", middleware.AdminOnly)
	inviteGroup.POST("", handlers.Invite.CreateInvite)
	inviteGroup.GET("", handlers.Invite.GetInvites)
	inviteGroup.DELETE("/:id", handlers.Invite.DeleteInvite)
}
