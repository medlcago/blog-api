package users

import (
	"blog-api/internal/errors"
	"blog-api/internal/logger"
	"blog-api/internal/photos"
	"blog-api/pkg/response"
	"context"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"
)

type IUserHandler interface {
	GetMe(ctx fiber.Ctx) error
	UploadAvatar(ctx fiber.Ctx) error
}

type UserHandler struct {
	photoService photos.IPhotoService
}

func NewUserHandler(photoService photos.IPhotoService) IUserHandler {
	return &UserHandler{
		photoService: photoService,
	}
}

func (h *UserHandler) GetMe(ctx fiber.Ctx) error {
	user := MustGetUser(ctx)

	return ctx.JSON(response.Response[*UserResponse]{
		OK:   true,
		Data: user,
	})
}

func (h *UserHandler) UploadAvatar(ctx fiber.Ctx) error {
	user := MustGetUser(ctx)
	requestID := requestid.FromContext(ctx)

	file, err := ctx.FormFile("file")
	if err != nil {
		return errors.ErrInvalidFile
	}

	res, err := h.photoService.UploadAvatar(
		context.WithValue(ctx, logger.RequestIDKey, requestID),
		user.UserID,
		file,
	)

	if err != nil {
		return err
	}

	return ctx.JSON(response.NewResponse(res))
}
