package routes

import (
	"blog-api/internal/auth"
	"blog-api/internal/middleware"

	"github.com/gofiber/fiber/v3"
)

func RegisterAuthRoutes(r fiber.Router, h auth.IAuthHandler, mw *middleware.Manager) {
	r.Post("/register", h.Register)
	r.Post("/login", h.Login)
	r.Post("/refresh-token", h.RefreshToken)
	r.Post("/change-password", mw.AuthMiddleware(), h.ChangePassword)
}
