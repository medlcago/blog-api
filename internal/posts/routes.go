package posts

import "github.com/gofiber/fiber/v3"

func RegisterRoutes(r fiber.Router, h IPostHandler) {
	r.Post("/", h.CreatePost)
	r.Get("/:id<int>", h.GetPost)
	r.Get("/", h.GetPosts)
}
