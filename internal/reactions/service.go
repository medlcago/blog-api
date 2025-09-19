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
	GetAvailableReactions(ctx context.Context) ([]*ReactionTypeResponse, error)
}

type reactionService struct {
	db     *database.DB
	logger *slog.Logger
}

func NewReactionService(db *database.DB, logger *slog.Logger) IReactionService {
	return &reactionService{
		db:     db,
		logger: logger,
	}
}

func (s *reactionService) setReaction(ctx context.Context, userID uint, input SetReactionInput) (*ReactionResponse, error) {
	log := logger.WithUserID(logger.FromCtx(ctx, s.logger), userID).With(
		slog.String("target_type", input.TargetType),
		slog.Uint64("target_id", uint64(input.TargetID)),
		slog.Uint64("reaction_type_id", uint64(input.ReactionID)),
	)

	if !slices.Contains(allowedTargets, input.TargetType) {
		log.Warn("invalid target type")
		return nil, errors.BadRequest("invalid target type")
	}

	db := s.db.WithContext(ctx)

	var reactType models.ReactionType
	if err := db.Where("id = ? AND is_active = ?", input.ReactionID, "true").First(&reactType).Error; err != nil {
		if goerrors.Is(err, database.ErrRecordNotFound) {
			log.Warn("invalid or inactive reaction")
			return nil, errors.BadRequest("invalid or inactive reaction")
		}
		return nil, err
	}

	var finalReaction *models.Reaction
	err := db.Transaction(func(tx *gorm.DB) error {
		var existing models.Reaction
		err := tx.Where("user_id = ? AND target_type = ? AND target_id = ?",
			userID, input.TargetType, input.TargetID,
		).First(&existing).Error

		switch {
		case err == nil && existing.ReactionTypeID == input.ReactionID:
			log.Info("removing existing reaction")
			if e := tx.Delete(&existing).Error; e != nil {
				log.Error("failed to delete reaction", logger.Err(e))
				return err
			}
			finalReaction = nil

		case err == nil:
			log.Info("updating reaction", slog.Uint64("new_reaction_type_id", uint64(input.ReactionID)))
			existing.ReactionTypeID = input.ReactionID
			if e := tx.Save(&existing).Error; e != nil {
				log.Error("failed to update reaction", logger.Err(e))
				return err
			}
			finalReaction = &existing

		case goerrors.Is(err, gorm.ErrRecordNotFound):
			log.Info("creating new reaction")
			newReaction := models.Reaction{
				UserID:         userID,
				TargetID:       input.TargetID,
				TargetType:     input.TargetType,
				ReactionTypeID: input.ReactionID,
			}
			if e := tx.Create(&newReaction).Error; e != nil {
				log.Error("failed to create reaction", logger.Err(e))
				return e
			}
			finalReaction = &newReaction

		default:
			log.Error("failed to fetch reaction", logger.Err(err))
			return err
		}
		return nil
	})

	aggReact, err := GetReactionsAggregate(db, input.TargetType, []uint{input.TargetID})
	if err != nil {
		log.Error("failed to aggregate reactions", logger.Err(err))
		return nil, err
	}
	reactions := aggReact[input.TargetID]

	reactResponse := &ReactionResponse{
		UserID:     userID,
		TargetID:   input.TargetID,
		TargetType: input.TargetType,
		Reactions:  reactions,
	}

	if finalReaction == nil { // reaction removed
		log.Info("reaction removed successfully")
		return reactResponse, nil
	}

	reactResponse.UserReaction = &models.UserReaction{
		TargetID: input.TargetID,
		Type:     reactType.Name,
		Icon:     reactType.Icon,
		IsActive: reactType.IsActive,
	}
	log.Info("reaction set successfully")
	return reactResponse, nil
}

func (s *reactionService) SetPostReaction(ctx context.Context, userID uint, input SetPostReactionInput) (*ReactionResponse, error) {
	return s.setReaction(ctx, userID,
		SetReactionInput{
			TargetType: TargetPost,
			TargetID:   input.PostID,
			ReactionID: input.ReactionID,
		},
	)
}

func (s *reactionService) GetAvailableReactions(ctx context.Context) ([]*ReactionTypeResponse, error) {
	db := s.db.WithContext(ctx)
	log := logger.FromCtx(ctx, s.logger)

	var availableReactionTypes []models.ReactionType
	if err := db.Find(&availableReactionTypes, "is_active = true").Error; err != nil {
		log.Error("failed to get available reaction types", logger.Err(err))
		return nil, err
	}

	return MapReactionTypesToResponse(availableReactionTypes), nil
}
