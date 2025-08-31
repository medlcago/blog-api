package middleware

import (
	"blog-api/internal/logger"
	"blog-api/internal/tokenmanager"
	"context"
	"log/slog"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"
)

func (m *Manager) OptionalAuthMiddleware() fiber.Handler {
	log := m.log.With(slog.String("component", "middleware/optional_auth"))

	return func(ctx fiber.Ctx) error {
		requestID := requestid.FromContext(ctx)
		reqLog := log.With(slog.String(string(logger.RequestIDKey), requestID))

		tokenHeader := fiber.GetReqHeader[string](ctx, "Authorization")
		if len(tokenHeader) < 7 || !strings.EqualFold(tokenHeader[:7], "bearer ") {
			reqLog.Info("missing bearer token; proceeding as guest")
			return ctx.Next()
		}

		claims, err := m.jwtService.ValidateToken(tokenHeader[7:])
		if err != nil {
			reqLog.Info("invalid token; proceeding as guest")
			return ctx.Next()
		}

		if claims.TokenType != tokenmanager.AccessToken {
			reqLog.Info("invalid token type; proceeding as guest", slog.String("token_type", claims.TokenType))
			return ctx.Next()
		}

		userID, err := strconv.ParseUint(claims.UserID, 10, 64)
		if err != nil {
			reqLog.Error("invalid user id; proceeding as guest", logger.Err(err))
			return ctx.Next()
		}

		user, err := m.userService.GetUserByID(
			context.WithValue(ctx, logger.RequestIDKey, requestID),
			uint(userID),
		)
		if err != nil {
			reqLog.Info("user lookup failed; proceeding as guest", logger.Err(err))
			return ctx.Next()
		}

		ctx.Locals("user", user)
		reqLog.Info("authenticated user", slog.Uint64("user_id", userID))
		return ctx.Next()
	}
}
