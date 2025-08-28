package posts

import (
	"blog-api/internal/users"
	"blog-api/pkg/errors"
	"blog-api/pkg/response"
	"context"

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
	user := users.MustGetUser(ctx)

	var input CreatePostInput

	if err := ctx.Bind().Body(&input); err != nil {
		return errors.ErrInvalidBody
	}

	post, err := h.postService.CreatePost(context.Background(), user.UserID, input)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusCreated).JSON(response.Response[*PostResponse]{
		OK:   true,
		Data: post,
	})
}

func (h *PostHandler) GetPost(ctx fiber.Ctx) error {
	postId := fiber.Params[uint](ctx, "id")
	post, err := h.postService.GetPost(context.Background(), postId)
	if err != nil {
		return err
	}

	return ctx.JSON(response.Response[*PostResponse]{
		OK:   true,
		Data: post,
	})
}

func (h *PostHandler) GetPosts(ctx fiber.Ctx) error {
	var params FilterParams
	if err := ctx.Bind().Query(&params); err != nil {
		return errors.ErrInvalidQuery
	}

	data, err := h.postService.GetPosts(context.Background(), params)
	if err != nil {
		return err
	}

	return ctx.JSON(response.NewPaginatedResponse(data.Total, data.Result))
}

func (h *PostHandler) UpdatePost(ctx fiber.Ctx) error {
	user := users.MustGetUser(ctx)

	var input CreatePostInput
	if err := ctx.Bind().Body(&input); err != nil {
		return errors.ErrInvalidBody
	}

	postID := fiber.Params[uint](ctx, "id")
	updatedPost, err := h.postService.UpdatePost(context.Background(), user.UserID, postID, input)
	if err != nil {
		return err
	}

	return ctx.JSON(response.Response[*PostResponse]{
		OK:   true,
		Data: updatedPost,
	})
}

func (h *PostHandler) DeletePost(ctx fiber.Ctx) error {
	user := users.MustGetUser(ctx)

	postID := fiber.Params[uint](ctx, "id")
	if err := h.postService.DeletePost(context.Background(), user.UserID, postID); err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
