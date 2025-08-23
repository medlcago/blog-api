package posts

import (
	"blog-api/internal/users"
	"blog-api/pkg/errors"
	"blog-api/pkg/response"

	"github.com/gofiber/fiber/v3"
)

type IPostHandler interface {
	CreatePost(ctx fiber.Ctx) error
	GetPost(ctx fiber.Ctx) error
	GetPosts(ctx fiber.Ctx) error
	UpdatePost(ctx fiber.Ctx) error
	DeletePost(ctx fiber.Ctx) error
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
	user := fiber.Locals[*users.UserResponse](ctx, "user")
	if user == nil {
		return errors.ErrUnauthorized
	}

	var input CreatePostInput

	if err := ctx.Bind().Body(&input); err != nil {
		return errors.ErrInvalidBody
	}

	res, err := h.postService.CreatePost(user.UserID, input)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusCreated).JSON(response.Response[*PostResponse]{
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

	return ctx.JSON(response.Response[*PostResponse]{
		OK:   true,
		Data: res,
	})
}

func (h *PostHandler) GetPosts(ctx fiber.Ctx) error {
	var filter FilterParams
	if err := ctx.Bind().Query(&filter); err != nil {
		return errors.ErrInvalidQuery
	}

	filter = DefaultFilterParams(filter)

	res, err := h.postService.GetPosts(&filter)
	if err != nil {
		return err
	}
	return ctx.JSON(response.Response[[]*PostResponse]{
		OK:   true,
		Data: res,
	})
}

func (h *PostHandler) UpdatePost(ctx fiber.Ctx) error {
	user := fiber.Locals[*users.UserResponse](ctx, "user")
	if user == nil {
		return errors.ErrUnauthorized
	}

	var input CreatePostInput
	if err := ctx.Bind().Body(&input); err != nil {
		return errors.ErrInvalidBody
	}

	postID := fiber.Params[uint](ctx, "id")
	newPost, err := h.postService.UpdatePost(user.UserID, postID, input)
	if err != nil {
		return err
	}

	return ctx.JSON(response.Response[*PostResponse]{
		OK:   true,
		Data: newPost,
	})
}

func (h *PostHandler) DeletePost(ctx fiber.Ctx) error {
	user := fiber.Locals[*users.UserResponse](ctx, "user")
	if user == nil {
		return errors.ErrUnauthorized
	}

	postID := fiber.Params[uint](ctx, "id")
	if err := h.postService.DeletePost(user.UserID, postID); err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
