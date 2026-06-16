package handlers

import (
	"net/http"

	"github.com/ComputerSocietyVITC/projects-portal-backend/internal/config"
	"github.com/ComputerSocietyVITC/projects-portal-backend/internal/models"
	"github.com/labstack/echo/v5"
	"go.uber.org/zap"
)

type UserHandler struct {
	DB *config.Database
}

func (h *UserHandler) GetUsers(c *echo.Context) error {
	var users []models.User
	if err := h.DB.Find(&users).Error; err != nil {
		zap.L().Error("failed to fetch users", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch users"})
	}
	return c.JSON(http.StatusOK, users)
}

func (h *UserHandler) GetUserByID(c *echo.Context) error {
	id := c.Param("id")
	var user models.User
	if err := h.DB.First(&user, id).Error; err != nil {
		zap.L().Error("user not found", zap.String("id", id), zap.Error(err))
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}
	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) CreateUser(c *echo.Context) error {
	var user models.User
	if err := c.Bind(&user); err != nil {
		zap.L().Error("failed to bind user data", zap.Error(err))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user data"})
	}
	if err := h.DB.Create(&user).Error; err != nil {
		zap.L().Error("failed to create user", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
	}
	return c.JSON(http.StatusCreated, user)
}
