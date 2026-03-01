package handlers

import (
	"errors"
	"net/http"

	"github.com/ComputerSocietyVITC/projects-portal-backend/internal/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"go.uber.org/zap"
)

type ProjectMemberHandler struct {
	service service.ProjectMemberService
	logger  *zap.Logger
}

func NewProjectMemberHandler(service service.ProjectMemberService, logger *zap.Logger) *ProjectMemberHandler {
	return &ProjectMemberHandler{
		service: service,
		logger:  logger,
	}
}

type AddMemberRequest struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
	Role   string    `json:"role" validate:"required,oneof=maintainer member viewer"`
}

func (h *ProjectMemberHandler) AddMember(c *echo.Context) error {
	projectID, err := echo.PathParam[uuid.UUID](c, "project_id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid project id"})
	}

	var req AddMemberRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
	}

	if err := h.service.AddUserToProject(c.Request().Context(), projectID, req.UserID, req.Role); err != nil {
		if errors.Is(err, service.ErrMemberAlreadyExists) {
			return c.JSON(http.StatusConflict, ErrorResponse{Error: "user is already a member"})
		}
		h.logger.Error("failed to add member", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to add member"})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "member added successfully"})
}

func (h *ProjectMemberHandler) RemoveMember(c *echo.Context) error {
	projectID, err := echo.PathParam[uuid.UUID](c, "project_id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid project id"})
	}

	userID, err := echo.PathParam[uuid.UUID](c, "user_id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid user id"})
	}

	if err := h.service.RemoveUserFromProject(c.Request().Context(), projectID, userID); err != nil {
		if errors.Is(err, service.ErrMemberNotFound) {
			return c.JSON(http.StatusNotFound, ErrorResponse{Error: "member not found"})
		}
		h.logger.Error("failed to remove member", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to remove member"})
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *ProjectMemberHandler) GetMembers(c *echo.Context) error {
	projectID, err := echo.PathParam[uuid.UUID](c, "project_id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid project id"})
	}

	members, err := h.service.GetProjectMembers(c.Request().Context(), projectID)
	if err != nil {
		h.logger.Error("failed to get members", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to get members"})
	}

	return c.JSON(http.StatusOK, members)
}
