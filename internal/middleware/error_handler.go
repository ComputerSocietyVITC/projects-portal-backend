package middleware

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v5"
	"go.uber.org/zap"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func ErrorHandler(logger *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			err := next(c)
			if err != nil {
				logger.Error("request error",
					zap.Error(err),
					zap.String("path", c.Request().URL.Path),
					zap.String("method", c.Request().Method),
				)

				if he, ok := err.(*echo.HTTPError); ok {
					return c.JSON(he.Code, ErrorResponse{Error: fmt.Sprintf("%v", he.Message)})
				}

				return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
			}
			return nil
		}
	}
}
