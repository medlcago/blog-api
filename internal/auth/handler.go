package auth

import (
	"blog-api/internal/logger"
	"blog-api/internal/users"
	"blog-api/pkg/response"
	"context"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"
)

type IAuthHandler interface {
	Register(ctx fiber.Ctx) error
	Login(ctx fiber.Ctx) error
	RefreshToken(ctx fiber.Ctx) error
	ChangePassword(ctx fiber.Ctx) error
}

type AuthHandler struct {
	authService IAuthService
}

func NewAuthHandler(authService IAuthService) IAuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(ctx fiber.Ctx) error {
	requestID := requestid.FromContext(ctx)

	var input RegisterUserInput
	if err := ctx.Bind().JSON(&input); err != nil {
		return err
	}

	token, err := h.authService.Register(
		context.WithValue(ctx, logger.RequestIDKey, requestID),
		input,
	)

	if err != nil {
		return err
	}

	return ctx.Status(201).JSON(response.Response[*TokenResponse]{
		OK:   true,
		Msg:  "User registered!",
		Data: token,
	})

}

func (h *AuthHandler) Login(ctx fiber.Ctx) error {
	requestID := requestid.FromContext(ctx)

	var input LoginUserInput

	if err := ctx.Bind().JSON(&input); err != nil {
		return err
	}

	token, err := h.authService.Login(
		context.WithValue(ctx, logger.RequestIDKey, requestID),
		input,
	)

	if err != nil {
		return err
	}

	return ctx.JSON(response.Response[*TokenResponse]{
		OK:   true,
		Msg:  "User logged in!",
		Data: token,
	})
}

func (h *AuthHandler) RefreshToken(ctx fiber.Ctx) error {
	requestID := requestid.FromContext(ctx)

	var input RefreshTokenInput

	if err := ctx.Bind().JSON(&input); err != nil {
		return err
	}

	res, err := h.authService.RefreshToken(
		context.WithValue(ctx, logger.RequestIDKey, requestID),
		input,
	)

	if err != nil {
		return err
	}

	return ctx.JSON(response.Response[*TokenResponse]{
		OK:   true,
		Msg:  "Token refreshed!",
		Data: res,
	})
}

func (h *AuthHandler) ChangePassword(ctx fiber.Ctx) error {
	requestID := requestid.FromContext(ctx)

	user := users.MustGetUser(ctx)

	var input ChangePasswordInput

	if err := ctx.Bind().Form(&input); err != nil {
		return err
	}

	err := h.authService.ChangePassword(context.WithValue(ctx, logger.RequestIDKey, requestID), user.UserID, input)
	if err != nil {
		return err
	}

	return ctx.SendString("OK")
}
