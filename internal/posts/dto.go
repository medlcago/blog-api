package posts

import (
	"blog-api/internal/users"
	"time"
)

type PostEntityInput struct {
	Offset int     `json:"offset" validate:"min=0"`
	Length int     `json:"length" validate:"min=1"`
	Type   string  `json:"type" validate:"required,oneof=bold italic spoiler link underline"`
	URL    *string `json:"url" validate:"omitempty,http_url,max=500"`
}

type CreatePostInput struct {
	Title   string `json:"title" validate:"required,min=1,max=255"`
	Content string `json:"content" validate:"required,min=1"`

	Entities []PostEntityInput `json:"entities" validate:"omitempty,dive"`
}

type PostResponse struct {
	ID        uint      `json:"id"`
	AuthorID  uint      `json:"author_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`

	Author   *users.UserResponse `json:"author,omitempty"`
	Entities []*PostEntityInput  `json:"entities"`
}
