package handlers

import (
	"net/http"

	"github.com/ComputerSocietyVITC/projects-portal-backend/internal/models"
	"github.com/ComputerSocietyVITC/projects-portal-backend/internal/service"
	"github.com/labstack/echo/v5"
	"go.uber.org/zap"
)

type ProjectHandler struct {
	service service.ProjectService
	logger  *zap.Logger
}

func NewProjectHandler(service service.ProjectService, logger *zap.Logger) *ProjectHandler {
	return &ProjectHandler{
		service: service,
		logger:  logger,
	}
}

func (h *ProjectHandler) GetAllProjects(c *echo.Context) error {
	projects, err := h.service.GetAllProjects(c.Request().Context())
	if err != nil {
		h.logger.Error("failed to get projects", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to retrieve projects"})
	}
	return c.JSON(http.StatusOK, projects)
}

func (h *ProjectHandler) GetProjectByID(c *echo.Context) error {
	id, err := echo.PathParam[uint](c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid project id"})
	}
	project, err := h.service.GetProjectByID(c.Request().Context(), id)
	if err != nil {
		h.logger.Error("failed to get project", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to retrieve project"})
	}
	return c.JSON(http.StatusOK, project)
}

func (h *ProjectHandler) CreateProject(c *echo.Context) error {
	var project models.Project
	if err := c.Bind(&project); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}
	createdProject, err := h.service.CreateProject(c.Request().Context(), &project)
	if err != nil {
		h.logger.Error("failed to create project", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create project"})
	}
	return c.JSON(http.StatusCreated, createdProject)
}

func (h *ProjectHandler) UpdateProject(c *echo.Context) error {
	id, err := echo.PathParam[uint](c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid project id"})
	}
	var updates map[string]interface{}
	if err := c.Bind(&updates); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}
	updatedProject, err := h.service.UpdateProject(c.Request().Context(), id, updates)
	if err != nil {
		h.logger.Error("failed to update project", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to update project"})
	}
	return c.JSON(http.StatusOK, updatedProject)
}

func (h *ProjectHandler) DeleteProject(c *echo.Context) error {
	id, err := echo.PathParam[uint](c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid project id"})
	}
	if err := h.service.DeleteProject(c.Request().Context(), id); err != nil {
		h.logger.Error("failed to delete project", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to delete project"})
	}
	return c.NoContent(http.StatusNoContent)
}
