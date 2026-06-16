package handlers

import (
	"net/http"

	"github.com/ComputerSocietyVITC/projects-portal-backend/internal/config"
	"github.com/ComputerSocietyVITC/projects-portal-backend/internal/models"
	"github.com/labstack/echo/v5"
	"go.uber.org/zap"
)

type ProjectMemberHandler struct {
	DB *config.Database
}

func (h *ProjectMemberHandler) AddMember(c *echo.Context) error {
	var member models.ProjectMember
	if err := c.Bind(&member); err != nil {
		zap.L().Error("failed to bind project member data", zap.Error(err))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid project member data"})
	}
	if err := h.DB.Create(&member).Error; err != nil {
		zap.L().Error("failed to add project member", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to add project member"})
	}
	return c.JSON(http.StatusCreated, member)
}

func (h *ProjectMemberHandler) RemoveMember(c *echo.Context) error {
	id := c.Param("id")
	if err := h.DB.Delete(&models.ProjectMember{}, id).Error; err != nil {
		zap.L().Error("failed to remove project member", zap.String("id", id), zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to remove project member"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Project member removed successfully"})
}
