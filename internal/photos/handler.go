package photos

import (
	"blog-api/internal/errors"
	"blog-api/internal/logger"
	"blog-api/internal/users"
	"blog-api/pkg/response"
	"context"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"
)

type IPhotoHandler interface {
	UploadAvatar(ctx fiber.Ctx) error
}

type PhotoHandler struct {
	photoService IPhotoService
}

func NewPhotoHandler(photoService IPhotoService) IPhotoHandler {
	return &PhotoHandler{photoService: photoService}

}

func (h *PhotoHandler) UploadAvatar(ctx fiber.Ctx) error {
	user := users.MustGetUser(ctx)
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
