package users

import (
	"blog-api/internal/photos"
	"blog-api/pkg/errors"
	"blog-api/pkg/response"
	"context"

	"github.com/gofiber/fiber/v3"
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

func (u *UserHandler) GetMe(ctx fiber.Ctx) error {
	user := MustGetUser(ctx)

	return ctx.JSON(response.Response[*UserResponse]{
		OK:   true,
		Data: user,
	})
}

func (u *UserHandler) UploadAvatar(ctx fiber.Ctx) error {
	user := MustGetUser(ctx)

	file, err := ctx.FormFile("file")
	if err != nil {
		return errors.New(400, "invalid file")
	}

	res, err := u.photoService.UploadAvatar(context.Background(), user.UserID, file)
	if err != nil {
		return err
	}

	return ctx.JSON(response.NewResponse(res))
}
