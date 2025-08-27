package middleware

import (
	"blog-api/internal/jwtmanager"
	"blog-api/pkg/errors"
	"context"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
)

func (m *Manager) AuthMiddleware() fiber.Handler {
	return func(ctx fiber.Ctx) error {
		tokenHeader := fiber.GetReqHeader[string](ctx, "Authorization")
		if len(tokenHeader) < 7 || !strings.EqualFold(tokenHeader[:7], "bearer ") {
			return errors.ErrMissingToken
		}

		claims, err := m.jwtManager.ValidateToken(tokenHeader[7:])
		if err != nil {
			return errors.ErrUnauthorized
		}

		if claims.TokenType != jwtmanager.AccessToken {
			return errors.ErrUnauthorized
		}

		userID, err := strconv.ParseUint(claims.UserID, 10, 64)
		if err != nil {
			return errors.ErrUnauthorized
		}

		user, err := m.userService.GetUserByID(context.Background(), uint(userID))
		if err != nil {
			return errors.ErrUnauthorized
		}

		ctx.Locals("user", user)
		return ctx.Next()
	}
}
