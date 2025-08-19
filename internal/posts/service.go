package posts

import (
	"blog-api/internal/database"
	"blog-api/internal/models"
	"blog-api/pkg/errors"
	goerrors "errors"

	"gorm.io/gorm"
)

type IPostService interface {
	CreatePost(input CreatePostInput) (PostResponse, error)
	GetPost(postID uint) (PostResponse, error)
	GetPosts() ([]PostResponse, error)
}

type PostService struct {
}

func NewPostService() IPostService {
	return &PostService{}
}

func (p *PostService) CreatePost(input CreatePostInput) (PostResponse, error) {
	instance := database.GetDb()

	post := models.Post{
		Title:   input.Title,
		Content: input.Content,
	}

	err := instance.Create(&post).Error
	if err != nil {
		return PostResponse{}, err
	}
	return PostResponse{
		ID:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		CreatedAt: post.CreatedAt,
	}, nil
}

func (p *PostService) GetPost(postID uint) (PostResponse, error) {
	instance := database.GetDb()

	var post models.Post

	err := instance.First(&post, postID).Error
	if err != nil {
		if goerrors.Is(err, gorm.ErrRecordNotFound) {
			return PostResponse{}, errors.New(404, "post not found")
		}
		return PostResponse{}, err
	}
	return PostResponse{
		ID:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		CreatedAt: post.CreatedAt,
	}, nil
}

func (p *PostService) GetPosts() ([]PostResponse, error) {
	instance := database.GetDb()

	var posts []models.Post
	output := make([]PostResponse, 0)

	err := instance.Find(&posts).Error
	if err != nil {
		return output, err
	}

	for _, post := range posts {
		output = append(output, PostResponse{
			ID:        post.ID,
			Title:     post.Title,
			Content:   post.Content,
			CreatedAt: post.CreatedAt,
		})
	}

	return output, nil
}
