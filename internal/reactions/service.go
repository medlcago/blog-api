package reactions

import (
	"blog-api/internal/database"
	"blog-api/internal/errors"
	"blog-api/internal/logger"
	"blog-api/internal/models"
	"context"
	goerrors "errors"
	"log/slog"
	"slices"

	"gorm.io/gorm"
)

const (
	TargetPost = "posts"
)

var (
	allowedTargets = []string{TargetPost}
)

type IReactionService interface {
	SetPostReaction(ctx context.Context, userID uint, input SetPostReactionInput) (*ReactionResponse, error)
}

type ReactionService struct {
	db     *database.DB
	logger *slog.Logger
}

func NewReactionService(db *database.DB, logger *slog.Logger) IReactionService {
	return &ReactionService{
		db:     db,
		logger: logger,
	}
}

func (s *ReactionService) setReaction(ctx context.Context, userID uint, input SetReactionInput) (*ReactionResponse, error) {
	log := logger.FromCtx(ctx, s.logger).With(
		slog.Uint64("user_id", uint64(userID)),
		slog.String("target_type", input.TargetType),
		slog.Uint64("target_id", uint64(input.TargetID)),
		slog.String("reaction_type", string(input.ReactionType)),
	)

	if !slices.Contains(allowedTargets, input.TargetType) {
		log.Warn("invalid target type")
		return nil, errors.New(400, "invalid target type")
	}

	db := s.db.Get().WithContext(ctx)

	var finalReaction *models.Reaction

	err := db.Transaction(func(tx *gorm.DB) error {
		var existing models.Reaction
		err := tx.Where("user_id = ? AND target_type = ? AND target_id = ?",
			userID, input.TargetType, input.TargetID,
		).First(&existing).Error

		if err == nil {
			if existing.Type == input.ReactionType {
				// delete if the same type
				log.Info("removing existing reaction")
				if err := tx.Delete(&existing).Error; err != nil {
					log.Error("failed to delete reaction", logger.Err(err))
					return err
				}
				finalReaction = nil
				return nil
			} else {
				log.Info("updating reaction type",
					slog.String("old_type", string(existing.Type)),
					slog.String("new_type", string(input.ReactionType)),
				)

				existing.Type = input.ReactionType
				if err := tx.Save(&existing).Error; err != nil {
					log.Error("failed to update reaction", logger.Err(err))
					return err
				}
				finalReaction = &existing
				return nil
			}
		} else if !goerrors.Is(err, gorm.ErrRecordNotFound) {
			log.Error("failed to fetch existing reaction", logger.Err(err))
			return err
		}

		log.Info("creating new reaction")
		newReaction := models.Reaction{
			UserID:     userID,
			Type:       input.ReactionType,
			TargetID:   input.TargetID,
			TargetType: input.TargetType,
		}
		if err := tx.Create(&newReaction).Error; err != nil {
			log.Error("failed to create reaction", logger.Err(err))
			return err
		}
		finalReaction = &newReaction
		return nil
	})
	if err != nil {
		return nil, err
	}

	aggReact, err := GetReactionsAggregate(db, TargetPost, []uint{input.TargetID})
	if err != nil {
		log.Error("failed to aggregate reactions", logger.Err(err))
		return nil, err
	}
	statistics := aggReact[input.TargetID]

	reactResponse := &ReactionResponse{
		UserID:     userID,
		TargetID:   input.TargetID,
		TargetType: input.TargetType,
		Statistics: statistics,
	}

	if finalReaction == nil { // reaction removed
		log.Info("reaction removed successfully")
		return reactResponse, nil
	}

	reactResponse.ReactionType = &finalReaction.Type
	log.Info("reaction set successfully",
		slog.String("reaction_type", string(finalReaction.Type)),
	)
	return reactResponse, nil
}

func (s *ReactionService) SetPostReaction(ctx context.Context, userID uint, input SetPostReactionInput) (*ReactionResponse, error) {
	return s.setReaction(ctx, userID,
		SetReactionInput{
			TargetType:   TargetPost,
			TargetID:     input.PostID,
			ReactionType: input.ReactionType,
		},
	)
}
