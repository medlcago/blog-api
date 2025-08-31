package posts

import (
	"blog-api/internal/models"
	"blog-api/internal/users"
)

func MapPostsToResponse(posts []models.Post) []*PostResponse {
	output := make([]*PostResponse, len(posts))
	for i, post := range posts {
		output[i] = MapPostToResponse(post)
	}
	return output
}

func MapPostToResponse(post models.Post) *PostResponse {
	if post.ID == 0 {
		return nil
	}

	author := users.MapUserToResponse(post.Author)
	return &PostResponse{
		ID:           post.ID,
		AuthorID:     post.AuthorID,
		Title:        post.Title,
		Content:      post.Content,
		CreatedAt:    post.CreatedAt,
		Author:       author,
		Entities:     MapEntitiesToResponse(post.Entities),
		UserReaction: post.UserReaction,
		Reactions:    post.Reactions,
	}
}

func MapEntitiesToResponse(entities []models.PostEntity) []*PostEntityInput {
	output := make([]*PostEntityInput, len(entities))
	for i, entity := range entities {
		output[i] = MapEntityToResponse(entity)
	}
	return output
}

func MapEntityToResponse(entity models.PostEntity) *PostEntityInput {
	if entity.ID == 0 {
		return nil
	}

	return &PostEntityInput{
		Offset: entity.Offset,
		Length: entity.Length,
		Type:   entity.Type,
		URL:    entity.URL,
	}
}

func MapInputToPostEntity(input PostEntityInput) models.PostEntity {
	return models.PostEntity{
		Offset: input.Offset,
		Length: input.Length,
		Type:   input.Type,
		URL:    input.URL,
	}
}

func MapInputsToPostEntity(inputs []PostEntityInput) []models.PostEntity {
	entities := make([]models.PostEntity, len(inputs))
	for i, entity := range inputs {
		entities[i] = MapInputToPostEntity(entity)
	}
	return entities
}
