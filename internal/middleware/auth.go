package middleware

import (
	"blog-api/pkg/errors"
	"blog-api/pkg/jwtmanager"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

func (m *Manager) AuthMiddleware() fiber.Handler {
	return func(ctx fiber.Ctx) error {
		token := fiber.GetReqHeader[string](ctx, "Authorization")
		if token == "" {
			return errors.ErrMissingToken
		}

		claims, err := m.jwtManager.ValidateToken(token)
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

		user, err := m.userService.GetUserByID(uint(userID))
		if err != nil {
			return errors.ErrUnauthorized
		}

		ctx.Locals("user", user)
		return ctx.Next()
	}
}
