package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/ComputerSocietyVITC/projects-portal-backend/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

func jwtKeyFunc(token *jwt.Token) (interface{}, error) {
	return []byte(os.Getenv("JWT_SECRET")), nil
}

func AdminOnly(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		authorizationHeader := c.Request().Header.Get("authorization")
		accessToken := (strings.Split(authorizationHeader, " "))

		if len(accessToken) != 2 {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid access token"})
		}
		parsedAccessToken, _ := jwt.ParseWithClaims(accessToken[1], &models.JwtUserAccessToken{}, jwtKeyFunc)

		claim, ok := parsedAccessToken.Claims.(*models.JwtUserAccessToken)
		if !ok || !parsedAccessToken.Valid {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid access token"})
		}
		if claim.Role != "admin" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "insufficient role permissions", "role": claim.Role})
		}

		c.Set("adminUUID", claim.UserID)
		return next(c)
	}
}

func LoggedIn(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		authorizationHeader := c.Request().Header.Get("authorization")
		accessToken := (strings.Split(authorizationHeader, " "))

		if len(accessToken) != 2 {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid access token"})
		}
		parsedAccessToken, _ := jwt.ParseWithClaims(accessToken[1], &models.JwtUserAccessToken{}, jwtKeyFunc)

		claim, ok := parsedAccessToken.Claims.(*models.JwtUserAccessToken)
		if !ok || !parsedAccessToken.Valid {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid access token"})
		}

		c.Set("userUUID", claim.UserID)
		return next(c)
	}
}
