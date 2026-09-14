package handlers

import (
	"net/http"

	"github.com/ComputerSocietyVITC/projects-portal-backend/internal/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"go.uber.org/zap"
)

type InviteHandler struct {
	Service *service.InviteService
	Logger  *zap.Logger
}

func NewInviteHandler(service *service.InviteService, logger *zap.Logger) *InviteHandler {
	return &InviteHandler{
		Service: service,
		Logger:  logger,
	}
}

type CreateInviteBody struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

func (h *InviteHandler) CreateInvite(c *echo.Context) error {
	var body CreateInviteBody
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "bad body parameters"})
	}

	adminUUID := uuid.MustParse(c.Get("adminUUID").(string))

	invite, err := h.Service.CreateInvite(body.Email, body.Role, adminUUID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"id":         invite.ID,
		"token":      invite.Token,
		"role":       invite.Role,
		"expires_at": invite.ExpiresAt,
	})
}

func (h *InviteHandler) GetInvites(c *echo.Context) error {
	invites, err := h.Service.GetInvites()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"invites": invites})
}

func (h *InviteHandler) DeleteInvite(c *echo.Context) error {
	inviteID := c.Param("id")
	if err := h.Service.DeleteInvite(uuid.MustParse(inviteID)); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusNoContent, nil)
}
