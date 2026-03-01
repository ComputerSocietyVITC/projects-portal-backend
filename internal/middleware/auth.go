package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

type JWTClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Roles  []string  `json:"roles"`
	jwt.RegisteredClaims
}

func AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "missing authorization header",
				})
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "invalid authorization format, use: Bearer <token>",
				})
			}

			token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(getJWTSecret()), nil
			})

			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "invalid or expired token",
				})
			}

			claims, ok := token.Claims.(*JWTClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "invalid token claims",
				})
			}

			c.Set("user_id", claims.UserID)
			c.Set("email", claims.Email)
			c.Set("roles", claims.Roles)

			return next(c)
		}
	}
}

func RequireRole(allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			userRoles, ok := c.Get("roles").([]string)
			if !ok {
				return c.JSON(http.StatusForbidden, map[string]string{
					"error": "no roles found in token",
				})
			}

			for _, userRole := range userRoles {
				for _, allowedRole := range allowedRoles {
					if userRole == allowedRole {
						return next(c)
					}
				}
			}

			return c.JSON(http.StatusForbidden, map[string]string{
				"error": "insufficient permissions",
			})
		}
	}
}

func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-secret-key-change-in-production"
	}
	return secret
}
