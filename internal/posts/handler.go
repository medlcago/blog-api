package posts

import (
	"blog-api/internal/models"
	"blog-api/pkg/errors"
	"blog-api/pkg/response"

	"github.com/gofiber/fiber/v3"
)

type IPostHandler interface {
	CreatePost(ctx fiber.Ctx) error
	GetPost(ctx fiber.Ctx) error
	GetPosts(ctx fiber.Ctx) error
}

type PostHandler struct {
	postService IPostService
}

func NewPostHandler(postService IPostService) IPostHandler {
	return &PostHandler{
		postService: postService,
	}
}

func (h *PostHandler) CreatePost(ctx fiber.Ctx) error {
	user := fiber.Locals[models.User](ctx, "user")
	if user.ID == 0 {
		return errors.ErrUnauthorized
	}

	var input CreatePostInput

	if err := ctx.Bind().Body(&input); err != nil {
		return err
	}

	res, err := h.postService.CreatePost(user.ID, input)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusCreated).JSON(response.Response[PostResponse]{
		OK:   true,
		Data: res,
	})
}

func (h *PostHandler) GetPost(ctx fiber.Ctx) error {
	postId := fiber.Params[uint](ctx, "id")
	res, err := h.postService.GetPost(postId)
	if err != nil {
		return err
	}

	return ctx.JSON(response.Response[PostResponse]{
		OK:   true,
		Data: res,
	})
}

func (h *PostHandler) GetPosts(ctx fiber.Ctx) error {
	res, err := h.postService.GetPosts()
	if err != nil {
		return err
	}
	return ctx.JSON(response.Response[[]PostResponse]{
		OK:   true,
		Data: res,
	})
}
