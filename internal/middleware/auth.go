package middleware

import (
	"blog-api/internal/errors"
	"blog-api/internal/logger"
	"blog-api/internal/tokenmanager"
	"context"
	"log/slog"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"
)

func (m *Manager) AuthMiddleware() fiber.Handler {
	log := m.log.With(slog.String("component", "middleware/auth"))

	return func(ctx fiber.Ctx) error {
		requestID := requestid.FromContext(ctx)
		log = log.With(
			slog.String(string(logger.RequestIDKey), requestID),
		)

		tokenHeader := fiber.GetReqHeader[string](ctx, "Authorization")
		if len(tokenHeader) < 7 || !strings.EqualFold(tokenHeader[:7], "bearer ") {
			log.Info("missing bearer token")
			return errors.ErrMissingToken
		}

		claims, err := m.jwtService.ValidateToken(tokenHeader[7:])
		if err != nil {
			log.Info("invalid token")
			return errors.ErrUnauthorized
		}

		if claims.TokenType != tokenmanager.AccessToken {
			log.Info("invalid token type", slog.String("token_type", claims.TokenType))
			return errors.ErrUnauthorized
		}

		userID, err := strconv.ParseUint(claims.UserID, 10, 64)
		if err != nil {
			log.Error("invalid user id", slog.Any("error", err))
			return errors.ErrUnauthorized
		}

		user, err := m.userService.GetUserByID(
			context.WithValue(ctx, logger.RequestIDKey, requestID),
			uint(userID),
		)

		if err != nil {
			return errors.ErrUnauthorized
		}

		ctx.Locals("user", user)
		return ctx.Next()
	}
}
