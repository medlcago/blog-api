package posts

import (
	"blog-api/internal/models"
	"blog-api/internal/users"
)

func MapPostsToResponse(posts []models.Post) []PostResponse {
	output := make([]PostResponse, len(posts))
	for i, post := range posts {
		output[i] = MapPostToResponse(post)
	}
	return output
}

func MapPostToResponse(post models.Post) PostResponse {
	author := users.MapUserToResponse(post.Author)
	return PostResponse{
		ID:        post.ID,
		AuthorID:  post.AuthorID,
		Title:     post.Title,
		Content:   post.Content,
		CreatedAt: post.CreatedAt,
		Author:    &author,
	}
}
