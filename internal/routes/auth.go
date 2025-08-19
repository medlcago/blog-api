package routes

import (
	"blog-api/internal/auth"

	"github.com/gofiber/fiber/v3"
)

func RegisterAuthRoutes(r fiber.Router, h auth.IAuthHandler) {
	r.Post("/register", h.Register)
	r.Post("/login", h.Login)
}
