package handlers

import (
	"errors"
	"net/http"

	"github.com/ComputerSocietyVITC/projects-portal-backend/internal/service"
	"github.com/labstack/echo/v5"
	"go.uber.org/zap"
)

type AuthHandler struct {
	service service.AuthService
	logger  *zap.Logger
}

func NewAuthHandler(service service.AuthService, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		service: service,
		logger:  logger,
	}
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Name     string `json:"name" validate:"required"`
	Token    string `json:"token" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (h *AuthHandler) Login(c *echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
	}

	token, err := h.service.Login(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			return c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid email or password"})
		}
		h.logger.Error("login failed", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "login failed"})
	}

	return c.JSON(http.StatusOK, LoginResponse{Token: token})
}

func (h *AuthHandler) Register(c *echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
	}

	user, err := h.service.Register(c.Request().Context(), req.Email, req.Password, req.Name, req.Token)
	if err != nil {
		switch err {
		case service.ErrInviteNotFound:
			return c.JSON(http.StatusForbidden, ErrorResponse{Error: "invalid invite token"})
		case service.ErrInviteUsed:
			return c.JSON(http.StatusForbidden, ErrorResponse{Error: "invite has already been used"})
		case service.ErrInviteExpired:
			return c.JSON(http.StatusForbidden, ErrorResponse{Error: "invite has expired"})
		case service.ErrInviteEmailMismatch:
			return c.JSON(http.StatusForbidden, ErrorResponse{Error: "email does not match invite"})
		}
		h.logger.Error("registration failed", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "registration failed"})
	}

	return c.JSON(http.StatusCreated, user)
}
