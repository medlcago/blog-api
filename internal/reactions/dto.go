package reactions

import "blog-api/internal/models"

type SetReactionInput struct {
	TargetType string `json:"target_type" validate:"required"`
	TargetID   uint   `json:"target_id" validate:"required"`
	ReactionID uint   `json:"reaction_id" validate:"required"`
}

type SetPostReactionInput struct {
	PostID     uint `json:"post_id" validate:"required"`
	ReactionID uint `json:"reaction_id" validate:"required"`
}

type ReactionResponse struct {
	UserID       uint                 `json:"user_id"`
	TargetType   string               `json:"target_type"`
	TargetID     uint                 `json:"target_id"`
	UserReaction *models.UserReaction `json:"user_reaction"`

	Reactions []models.ReactionStat `json:"reactions"`
}

type ReactionTypeResponse struct {
	ReactionID uint   `json:"reaction_id"`
	Type       string `json:"type"`
	Icon       string `json:"icon"`
	IsActive   bool   `json:"is_active"`
}
