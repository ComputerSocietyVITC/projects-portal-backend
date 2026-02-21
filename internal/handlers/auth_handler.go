package handlers

import (
	"net/http"

	"github.com/ComputerSocietyVITC/projects-portal-backend/internal/service"
	"github.com/labstack/echo/v5"
	"go.uber.org/zap"
)

type AuthHandler struct {
	Service *service.AuthService
}

func NewAuthHandler(service *service.AuthService, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		Service: service,
	}
}

type loginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(c *echo.Context) error {
	var body loginBody
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "bad body parameters"})
	}

	user, err := h.Service.Login(body.Email, body.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	refreshToken, err := h.Service.CreateRefreshToken(user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	accessToken, err := h.Service.CreateAccessToken(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"access_token": accessToken, "refresh_token": refreshToken})
}

type registerBody struct {
	Token    string `query:"invite_token"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

func (h *AuthHandler) Register(c *echo.Context) error {
	var body registerBody
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "bad body parameters"})
	}

	user, err := h.Service.Register(body.Email, body.Password, body.Name, body.Token)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	refreshToken, err := h.Service.CreateRefreshToken(user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	accessToken, err := h.Service.CreateAccessToken(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{"access_token": accessToken, "refresh_token": refreshToken})
}

type refreshBody struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *AuthHandler) Refresh(c *echo.Context) error {
	var body refreshBody
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "bad body parameters"})
	}

	user, err := h.Service.GetUserFromRefreshToken(body.RefreshToken)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	refreshToken, err := h.Service.CreateRefreshToken(user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	accessToken, err := h.Service.CreateAccessToken(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"access_token": accessToken, "refresh_token": refreshToken})
}

func (h *AuthHandler) Logout(c *echo.Context) error {
	// removes the refresh token from the redisDB and access token is shortlived so clearing from frontend is enough
	if err := h.Service.Logout(c.Get("userUUID").(string)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusNoContent, nil)
}
