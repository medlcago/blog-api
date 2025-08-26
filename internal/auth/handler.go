package auth

import (
	"blog-api/pkg/errors"
	"blog-api/pkg/response"
	"context"

	"github.com/gofiber/fiber/v3"
)

type IAuthHandler interface {
	Register(ctx fiber.Ctx) error
	Login(ctx fiber.Ctx) error
}

type AuthHandler struct {
	authService IAuthService
}

func NewAuthHandler(authService IAuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (a *AuthHandler) Register(ctx fiber.Ctx) error {
	var input RegisterUserInput

	if err := ctx.Bind().JSON(&input); err != nil {
		return errors.ErrInvalidBody
	}

	token, err := a.authService.Register(context.Background(), input)
	if err != nil {
		return err
	}

	return ctx.Status(201).JSON(response.Response[*TokenResponse]{
		OK:   true,
		Msg:  "User registered!",
		Data: token,
	})

}

func (a *AuthHandler) Login(ctx fiber.Ctx) error {
	var input LoginUserInput

	if err := ctx.Bind().JSON(&input); err != nil {
		return errors.ErrInvalidBody
	}

	token, err := a.authService.Login(context.Background(), input)
	if err != nil {
		return err
	}

	return ctx.JSON(response.Response[*TokenResponse]{
		OK:   true,
		Msg:  "User logged in!",
		Data: token,
	})
}
