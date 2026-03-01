package handlers

import (
	"errors"
	"net/http"

	"github.com/ComputerSocietyVITC/projects-portal-backend/internal/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"go.uber.org/zap"
)

type UserHandler struct {
	service service.UserService
	logger  *zap.Logger
}

func NewUserHandler(service service.UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		service: service,
		logger:  logger,
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type UpdateUserRequest struct {
	Name   string `json:"name,omitempty"`
	Status string `json:"status,omitempty"`
}

func (h *UserHandler) GetAllUsers(c *echo.Context) error {
	users, err := h.service.GetAllUsers(c.Request().Context())
	if err != nil {
		h.logger.Error("failed to get users", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to retrieve users"})
	}

	return c.JSON(http.StatusOK, users)
}

func (h *UserHandler) GetUserByID(c *echo.Context) error {
	id, err := echo.PathParam[uuid.UUID](c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid user id"})
	}

	// Members may only view their own profile; group_head can view anyone's
	callerID, _ := c.Get("user_id").(uuid.UUID)
	callerRoles, _ := c.Get("roles").([]string)
	isGroupHead := false
	for _, r := range callerRoles {
		if r == "group_head" {
			isGroupHead = true
			break
		}
	}
	if !isGroupHead && callerID != id {
		return c.JSON(http.StatusForbidden, ErrorResponse{Error: "insufficient permissions"})
	}

	user, err := h.service.GetUserByID(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, ErrorResponse{Error: "user not found"})
		}
		h.logger.Error("failed to get user", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to retrieve user"})
	}

	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(c *echo.Context) error {
	id, err := echo.PathParam[uuid.UUID](c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid user id"})
	}

	var req UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}

	user, err := h.service.UpdateUser(c.Request().Context(), id, updates)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, ErrorResponse{Error: "user not found"})
		}
		h.logger.Error("failed to update user", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to update user"})
	}

	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) DeleteUser(c *echo.Context) error {
	id, err := echo.PathParam[uuid.UUID](c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid user id"})
	}

	if err := h.service.DeleteUser(c.Request().Context(), id); err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, ErrorResponse{Error: "user not found"})
		}
		h.logger.Error("failed to delete user", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to delete user"})
	}

	return c.NoContent(http.StatusNoContent)
}
