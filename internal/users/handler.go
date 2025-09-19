package users

import (
	"blog-api/pkg/response"

	"github.com/gofiber/fiber/v3"
)

type IUserHandler interface {
	GetMe(ctx fiber.Ctx) error
}

type userHandler struct {
}

func NewUserHandler() IUserHandler {
	return &userHandler{}
}

func (h *userHandler) GetMe(ctx fiber.Ctx) error {
	user := MustGetUser(ctx)
	return ctx.JSON(response.NewResponse(user))
}
