package logger

import (
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func WithRequestLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			logger := log.With().
				Str("method", c.Request().Method).
				Str("path", c.Path()).
				Logger()
			ctx := logger.WithContext(c.Request().Context())
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}
