package posts

import (
	"blog-api/internal/database"
	"blog-api/internal/models"
	"blog-api/internal/users"
	"blog-api/pkg/errors"
	goerrors "errors"

	"gorm.io/gorm"
)

type IPostService interface {
	CreatePost(userID uint, input CreatePostInput) (PostResponse, error)
	GetPost(postID uint) (PostResponse, error)
	GetPosts() ([]PostResponse, error)
}

type PostService struct {
}

func NewPostService() IPostService {
	return &PostService{}
}

func (p *PostService) CreatePost(userID uint, input CreatePostInput) (PostResponse, error) {
	instance := database.GetDb()

	post := models.Post{
		Title:    input.Title,
		Content:  input.Content,
		AuthorID: userID,
	}

	err := instance.Create(&post).Error
	if err != nil {
		return PostResponse{}, err
	}
	return PostResponse{
		ID:        post.ID,
		AuthorID:  post.AuthorID,
		Title:     post.Title,
		Content:   post.Content,
		CreatedAt: post.CreatedAt,
	}, nil
}

func (p *PostService) GetPost(postID uint) (PostResponse, error) {
	instance := database.GetDb()

	var post models.Post

	err := instance.Preload("Author").First(&post, postID).Error
	if err != nil {
		if goerrors.Is(err, gorm.ErrRecordNotFound) {
			return PostResponse{}, errors.ErrNotFound
		}
		return PostResponse{}, err
	}
	return PostResponse{
		ID:        post.ID,
		AuthorID:  post.AuthorID,
		Title:     post.Title,
		Content:   post.Content,
		CreatedAt: post.CreatedAt,
		Author: &users.UserResponse{
			UserID:   post.Author.ID,
			Username: post.Author.Username,
			Email:    post.Author.Email.String,
			Deleted:  post.Author.DeletedAt.Valid,
		},
	}, nil
}

func (p *PostService) GetPosts() ([]PostResponse, error) {
	instance := database.GetDb()

	var posts []models.Post

	err := instance.Preload("Author").Find(&posts).Error
	if err != nil {
		return nil, err
	}

	output := make([]PostResponse, 0)
	for _, post := range posts {
		output = append(output, PostResponse{
			ID:        post.ID,
			AuthorID:  post.AuthorID,
			Title:     post.Title,
			Content:   post.Content,
			CreatedAt: post.CreatedAt,
			Author: &users.UserResponse{
				UserID:   post.Author.ID,
				Username: post.Author.Username,
				Email:    post.Author.Email.String,
				Deleted:  post.Author.DeletedAt.Valid,
			},
		})
	}

	return output, nil
}
