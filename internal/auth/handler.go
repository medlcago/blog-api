package auth

import (
	"blog-api/internal/logger"
	"blog-api/pkg/errors"
	"blog-api/pkg/response"
	"context"
	"log/slog"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"
)

type IAuthHandler interface {
	Register(ctx fiber.Ctx) error
	Login(ctx fiber.Ctx) error
	RefreshToken(ctx fiber.Ctx) error
}

type AuthHandler struct {
	authService IAuthService
	logger      *slog.Logger
}

func NewAuthHandler(authService IAuthService, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
	}
}

func (a *AuthHandler) Register(ctx fiber.Ctx) error {
	requestID := requestid.FromContext(ctx)
	log := logger.WithRequestID(
		a.logger,
		requestID,
	)

	var input RegisterUserInput
	if err := ctx.Bind().JSON(&input); err != nil {
		log.Error("invalid request body",
			slog.String("path", ctx.Path()),
			slog.Any("error", err),
		)
		return errors.ErrInvalidBody
	}

	log.Info("register attempt",
		slog.String("username", input.Username),
	)

	token, err := a.authService.Register(
		context.WithValue(ctx, logger.RequestIDKey, requestID),
		input,
	)

	if err != nil {
		return err
	}

	log.Info("user registered successfully", slog.String("username", input.Username))

	return ctx.Status(201).JSON(response.Response[*TokenResponse]{
		OK:   true,
		Msg:  "User registered!",
		Data: token,
	})

}

func (a *AuthHandler) Login(ctx fiber.Ctx) error {
	requestID := requestid.FromContext(ctx)
	log := logger.WithRequestID(
		a.logger,
		requestID,
	)

	var input LoginUserInput

	if err := ctx.Bind().JSON(&input); err != nil {
		log.Error("invalid request body",
			slog.String("path", ctx.Path()),
			slog.Any("error", err),
		)
		return errors.ErrInvalidBody
	}

	log.Info("login attempt",
		slog.String("username", input.Username),
	)

	token, err := a.authService.Login(
		context.WithValue(ctx, logger.RequestIDKey, requestID),
		input,
	)

	if err != nil {
		return err
	}

	log.Info("user successfully logged in", slog.String("username", input.Username))

	return ctx.JSON(response.Response[*TokenResponse]{
		OK:   true,
		Msg:  "User logged in!",
		Data: token,
	})
}

func (a *AuthHandler) RefreshToken(ctx fiber.Ctx) error {
	requestID := requestid.FromContext(ctx)
	log := logger.WithRequestID(
		a.logger,
		requestID,
	)

	var input RefreshTokenInput

	if err := ctx.Bind().JSON(&input); err != nil {
		log.Error("invalid request body",
			slog.String("path", ctx.Path()),
			slog.Any("error", err),
		)
		return errors.ErrInvalidBody
	}

	res, err := a.authService.RefreshToken(
		context.WithValue(ctx, logger.RequestIDKey, requestID),
		input,
	)

	if err != nil {
		return err
	}

	log.Info("token successfully refreshed")

	return ctx.JSON(response.Response[*TokenResponse]{
		OK:   true,
		Msg:  "Token refreshed!",
		Data: res,
	})
}
