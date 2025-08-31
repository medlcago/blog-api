package routes

import (
	"blog-api/internal/middleware"
	"blog-api/internal/reactions"

	"github.com/gofiber/fiber/v3"
)

func RegisterReactionRoutes(r fiber.Router, h reactions.IReactionHandler, mw *middleware.Manager) {
	r.Post("/posts", mw.AuthMiddleware(), h.SetPostReaction)
}
