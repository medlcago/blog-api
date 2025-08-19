package posts

import (
	"blog-api/internal/users"
	"time"
)

type CreatePostInput struct {
	Title   string `json:"title" validate:"required,min=1,max=255"`
	Content string `json:"content" validate:"required,min=1"`
}

type PostResponse struct {
	ID        uint      `json:"id"`
	AuthorID  uint      `json:"author_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`

	Author *users.UserResponse `json:"author,omitempty"`
}
