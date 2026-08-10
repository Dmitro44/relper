package middlewares

import (
	"errors"
	"log/slog"
	"time"

	apperr "github.com/Dmitro44/relper/internal/errors"
	"github.com/gofiber/fiber/v3"
)

func RequestLogger() fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()
		err := c.Next()

		status := c.Response().StatusCode()
		if err != nil {
			if fiberErr, ok := errors.AsType[*fiber.Error](err); ok {
				status = fiberErr.Code
			}
		}

		level := slog.LevelDebug
		if status >= 500 {
			level = slog.LevelError
		} else if status >= 400 {
			level = slog.LevelWarn
		}

		attrs := []any{
			"method", c.Method(),
			"path", c.Path(),
			"status", status,
			"duration_ms", time.Since(start).Milliseconds(),
		}
		if err != nil {
			attrs = append(attrs, "error", err)
			if appErr, ok := errors.AsType[*apperr.AppError](err); ok {
				attrs = append(attrs, "code", appErr.Code)
			}
		}

		slog.Log(c.Context(), level, "request", attrs...)
		return err
	}
}
