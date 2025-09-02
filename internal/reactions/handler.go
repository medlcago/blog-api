package reactions

import (
	"blog-api/internal/logger"
	"blog-api/internal/users"
	"blog-api/pkg/response"
	"context"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"
)

type IReactionHandler interface {
	SetPostReaction(ctx fiber.Ctx) error
	GetAvailableReactions(ctx fiber.Ctx) error
}

type ReactionHandler struct {
	reactionService IReactionService
}

func NewReactionHandler(reactionService IReactionService) IReactionHandler {
	return &ReactionHandler{
		reactionService: reactionService,
	}
}

func (h *ReactionHandler) SetPostReaction(ctx fiber.Ctx) error {
	user := users.MustGetUser(ctx)

	requestID := requestid.FromContext(ctx)

	var input SetPostReactionInput

	if err := ctx.Bind().Body(&input); err != nil {
		return err
	}

	res, err := h.reactionService.SetPostReaction(
		context.WithValue(ctx, logger.RequestIDKey, requestID),
		user.UserID,
		input,
	)

	if err != nil {
		return err
	}

	return ctx.JSON(response.NewResponse(res))
}

func (h *ReactionHandler) GetAvailableReactions(ctx fiber.Ctx) error {
	requestID := requestid.FromContext(ctx)

	res, err := h.reactionService.GetAvailableReactions(context.WithValue(ctx, logger.RequestIDKey, requestID))
	if err != nil {
		return err
	}

	return ctx.JSON(response.NewResponse(res))
}
