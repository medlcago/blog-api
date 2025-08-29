package middleware

import (
	"blog-api/internal/logger"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"
)

func (m *Manager) LoggerMiddleware() fiber.Handler {
	log := m.log.With(slog.String("component", "middleware/logger"))

	return func(c fiber.Ctx) error {

		start := time.Now()

		err := c.Next()
		if err != nil {
			errHandler := c.App().ErrorHandler
			if err := errHandler(c, err); err != nil {
				_ = c.SendStatus(fiber.StatusInternalServerError)
			}
		}

		entry := log.With(
			slog.String("method", c.Method()),
			slog.String("path", c.Path()),
			slog.String("remote_addr", c.IP()),
			slog.String("user_agent", c.Get("User-Agent")),
			slog.String(string(logger.RequestIDKey), requestid.FromContext(c)),
		)

		entry.Info("request completed",
			slog.Int("bytes", len(c.Response().Body())),
			slog.Int("status_code", c.Response().StatusCode()),
			slog.String("duration", time.Since(start).String()),
		)

		return err
	}
}
