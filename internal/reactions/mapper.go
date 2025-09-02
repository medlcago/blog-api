package reactions

import "blog-api/internal/models"

func MapReactionTypesToResponse(reactionTypes []models.ReactionType) []*ReactionTypeResponse {
	output := make([]*ReactionTypeResponse, len(reactionTypes))
	for i, reactionType := range reactionTypes {
		output[i] = MapReactionTypeToResponse(reactionType)
	}
	return output
}

func MapReactionTypeToResponse(reactionType models.ReactionType) *ReactionTypeResponse {
	if reactionType.ID == 0 {
		return nil
	}

	return &ReactionTypeResponse{
		ReactionID: reactionType.ID,
		Type:       reactionType.Name,
		Icon:       reactionType.Icon,
		IsActive:   reactionType.IsActive,
	}
}
