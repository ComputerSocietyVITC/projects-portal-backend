package handlers

import (
	"net/http"

	"github.com/ComputerSocietyVITC/projects-portal-backend/internal/config"
	"github.com/labstack/echo/v5"
	"go.uber.org/zap"
)

type ProjectHandler struct {
	DB *config.Database
}

func (h *ProjectHandler) GetProjects(c echo.Context) error {
	var projects []Project
	if err := h.DB.Find(&projects).Error; err != nil {
		zap.L().Error("failed to fetch projects", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch projects"})
	}
	return c.JSON(http.StatusOK, projects)
}

func (h *ProjectHandler) GetProjectByID(c echo.Context) error {
	id := c.Param("id")
	var project Project
	if err := h.DB.First(&project, id).Error; err != nil {
		zap.L().Error("project not found", zap.String("id", id), zap.Error(err))
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Project not found"})
	}
	return c.JSON(http.StatusOK, project)
}

func (h *ProjectHandler) CreateProject(c echo.Context) error {
	var project Project
	if err := c.Bind(&project); err != nil {
		zap.L().Error("failed to bind project data", zap.Error(err))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid project data"})
	}
	if err := h.DB.Create(&project).Error; err != nil {
		zap.L().Error("failed to create project", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create project"})
	}
	return c.JSON(http.StatusCreated, project)
}

func (h *ProjectHandler) UpdateProject(c echo.Context) error {
	id := c.Param("id")
	var project Project
	if err := h.DB.First(&project, id).Error; err != nil {
		zap.L().Error("project not found", zap.String("id", id), zap.Error(err))
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Project not found"})
	}
	if err := c.Bind(&project); err != nil {
		zap.L().Error("failed to bind project data", zap.Error(err))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid project data"})
	}
	if err := h.DB.Save(&project).Error; err != nil {
		zap.L().Error("failed to update project", zap.String("id", id), zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update project"})
	}
	return c.JSON(http.StatusOK, project)
}

func (h *ProjectHandler) DeleteProject(c echo.Context) error {
	id := c.Param("id")
	var project Project
	if err := h.DB.First(&project, id).Error; err != nil {
		zap.L().Error("project not found", zap.String("id", id), zap.Error(err))
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Project not found"})
	}
	if err := h.DB.Delete(&project).Error; err != nil {
		zap.L().Error("failed to delete project", zap.String("id", id), zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete project"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Project deleted successfully"})
}
