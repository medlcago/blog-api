package routes

import (
	"blog-api/internal/middleware"
	"blog-api/internal/users"

	"github.com/gofiber/fiber/v3"
)

func RegisterUserRoutes(r fiber.Router, h users.IUserHandler, mw *middleware.Manager) {
	r.Get("/me", mw.AuthMiddleware(), h.GetMe)
}
