package reactions

import "blog-api/internal/models"

type SetReactionInput struct {
	TargetType   string              `json:"target_type" validate:"required"`
	TargetID     uint                `json:"target_id" validate:"required"`
	ReactionType models.ReactionType `json:"reaction_type" validate:"reaction"`
}

type SetPostReactionInput struct {
	PostID       uint                `json:"post_id" validate:"required"`
	ReactionType models.ReactionType `json:"reaction_type" validate:"reaction"`
}

type ReactionResponse struct {
	UserID       uint                 `json:"user_id"`
	TargetType   string               `json:"target_type"`
	TargetID     uint                 `json:"target_id"`
	ReactionType *models.ReactionType `json:"reaction_type"`

	Statistics map[string]int64 `json:"statistics"` // {"like": 10,  "love": 2}
}
