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

	Login2FA(ctx fiber.Ctx) error
	Enable2FA(ctx fiber.Ctx) error
	Verify2FA(ctx fiber.Ctx) error
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

	res, err := h.authService.Login(
		context.WithValue(ctx, logger.RequestIDKey, requestID),
		input,
	)

	if err != nil {
		return err
	}

	return ctx.JSON(response.NewResponse(res))
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

func (h *AuthHandler) Login2FA(ctx fiber.Ctx) error {
	requestID := requestid.FromContext(ctx)

	var input Login2FAInput

	if err := ctx.Bind().JSON(&input); err != nil {
		return err
	}

	res, err := h.authService.Login2FA(
		context.WithValue(ctx, logger.RequestIDKey, requestID),
		input,
	)
	if err != nil {
		return err
	}

	return ctx.JSON(response.NewResponse(res))
}

func (h *AuthHandler) Enable2FA(ctx fiber.Ctx) error {
	user := users.MustGetUser(ctx)

	requestID := requestid.FromContext(ctx)

	res, err := h.authService.Enable2FA(
		context.WithValue(ctx, logger.RequestIDKey, requestID),
		user.UserID,
	)
	if err != nil {
		return err
	}

	return ctx.JSON(response.NewResponse(res))
}

func (h *AuthHandler) Verify2FA(ctx fiber.Ctx) error {
	user := users.MustGetUser(ctx)

	requestID := requestid.FromContext(ctx)

	var input Verify2FAInput
	if err := ctx.Bind().JSON(&input); err != nil {
		return err
	}

	err := h.authService.Verify2FA(
		context.WithValue(ctx, logger.RequestIDKey, requestID),
		user.UserID,
		input,
	)
	if err != nil {
		return err
	}

	return ctx.SendString("OK")
}
