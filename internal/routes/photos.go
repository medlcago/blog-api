package routes

import (
	"blog-api/internal/middleware"
	"blog-api/internal/photos"

	"github.com/gofiber/fiber/v3"
)

func RegisterPhotoRoutes(r fiber.Router, h photos.IPhotoHandler, mw *middleware.Manager) {
	r.Post("/avatar", mw.AuthMiddleware(), h.UploadAvatar)
}
