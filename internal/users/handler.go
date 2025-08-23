package users

import (
	"blog-api/pkg/errors"
	"blog-api/pkg/response"

	"github.com/gofiber/fiber/v3"
)

type IUserHandler interface {
	GetMe(ctx fiber.Ctx) error
}

type UserHandler struct {
}

func NewUserHandler() IUserHandler {
	return &UserHandler{}
}

func (u *UserHandler) GetMe(ctx fiber.Ctx) error {
	user := fiber.Locals[*UserResponse](ctx, "user")
	if user == nil {
		return errors.ErrUnauthorized
	}

	return ctx.JSON(response.Response[*UserResponse]{
		OK:   true,
		Data: user,
	})
}
