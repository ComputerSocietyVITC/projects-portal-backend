package handlers

import (
	"net/http"

	"github.com/ComputerSocietyVITC/projects-portal-backend/internal/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"go.uber.org/zap"
)

type InviteHandler struct {
	service service.InviteService
	logger  *zap.Logger
}

func NewInviteHandler(service service.InviteService, logger *zap.Logger) *InviteHandler {
	return &InviteHandler{
		service: service,
		logger:  logger,
	}
}

type CreateInviteRequest struct {
	Email string `json:"email" validate:"required,email"`
	Role  string `json:"role" validate:"required"`
}

func (h *InviteHandler) CreateInvite(c *echo.Context) error {
	var req CreateInviteRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
	}

	createdBy, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
	}

	invite, err := h.service.CreateInvite(c.Request().Context(), req.Email, req.Role, createdBy)
	if err != nil {
		h.logger.Error("failed to create invite", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to create invite"})
	}

	return c.JSON(http.StatusCreated, invite)
}
