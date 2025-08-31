package posts

import (
	"blog-api/internal/errors"
	"blog-api/internal/logger"
	"blog-api/internal/users"
	"blog-api/pkg/response"
	"context"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"
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

	requestID := requestid.FromContext(ctx)

	var input CreatePostInput

	if err := ctx.Bind().Body(&input); err != nil {
		return err
	}

	post, err := h.postService.CreatePost(
		context.WithValue(ctx, logger.RequestIDKey, requestID), user.UserID,
		input,
	)

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
	requestID := requestid.FromContext(ctx)

	user := users.GetUser(ctx)
	var userID *uint
	if user != nil {
		userID = &user.UserID
	}

	post, err := h.postService.GetPost(
		context.WithValue(ctx, logger.RequestIDKey, requestID),
		postId,
		userID,
	)

	if err != nil {
		return err
	}

	return ctx.JSON(response.Response[*PostResponse]{
		OK:   true,
		Data: post,
	})
}

func (h *PostHandler) GetPosts(ctx fiber.Ctx) error {
	requestID := requestid.FromContext(ctx)

	user := users.GetUser(ctx)
	var userID *uint
	if user != nil {
		userID = &user.UserID
	}

	var params FilterParams
	if err := ctx.Bind().Query(&params); err != nil {
		return errors.ErrInvalidQuery
	}

	data, err := h.postService.GetPosts(
		context.WithValue(ctx, logger.RequestIDKey, requestID),
		params,
		userID,
	)

	if err != nil {
		return err
	}

	return ctx.JSON(response.NewPaginatedResponse(data.Total, data.Result))
}

func (h *PostHandler) UpdatePost(ctx fiber.Ctx) error {
	user := users.MustGetUser(ctx)
	postID := fiber.Params[uint](ctx, "id")

	requestID := requestid.FromContext(ctx)

	var input CreatePostInput
	if err := ctx.Bind().Body(&input); err != nil {
		return err
	}

	updatedPost, err := h.postService.UpdatePost(
		context.WithValue(ctx, logger.RequestIDKey, requestID),
		user.UserID,
		postID,
		input,
	)

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

	requestID := requestid.FromContext(ctx)

	err := h.postService.DeletePost(
		context.WithValue(ctx, logger.RequestIDKey, requestID),
		user.UserID,
		postID,
	)

	if err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
