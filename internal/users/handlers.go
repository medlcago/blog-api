package users

import (
	"blog-api/internal/models"
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
	user := fiber.Locals[models.User](ctx, "user")
	if user.ID == 0 {
		return errors.ErrUnauthorized
	}

	return ctx.JSON(response.Response[UserResponse]{
		OK: true,
		Data: UserResponse{
			UserID:   user.ID,
			Username: user.Username,
			Email:    user.Email.String,
			Deleted:  user.DeletedAt.Valid,
		},
	})
}
