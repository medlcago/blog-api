package routes

import (
	"blog-api/internal/middleware"
	"blog-api/internal/posts"

	"github.com/gofiber/fiber/v3"
)

func RegisterPostRoutes(r fiber.Router, h posts.IPostHandler, middlewareManager *middleware.Manager) {
	r.Post("/", middlewareManager.AuthMiddleware(), h.CreatePost)
	r.Get("/:id<int>", h.GetPost)
	r.Get("/", h.GetPosts)
}
