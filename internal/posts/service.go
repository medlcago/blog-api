package posts

import (
	"blog-api/internal/database"
	"blog-api/internal/models"
	"blog-api/pkg/errors"
	goerrors "errors"

	"gorm.io/gorm"
)

type IPostService interface {
	CreatePost(userID uint, input CreatePostInput) (PostResponse, error)
	GetPost(postID uint) (PostResponse, error)
	GetPosts(filter database.Filter) ([]PostResponse, error)
}

type PostService struct {
	db *database.DB
}

func NewPostService(db *database.DB) IPostService {
	return &PostService{
		db: db,
	}
}

func (p *PostService) CreatePost(userID uint, input CreatePostInput) (PostResponse, error) {
	db := p.db.Get()

	post := models.Post{
		Title:    input.Title,
		Content:  input.Content,
		AuthorID: userID,
	}

	err := db.Create(&post).Error
	if err != nil {
		return PostResponse{}, err
	}
	return MapPostToResponse(post), nil
}

func (p *PostService) GetPost(postID uint) (PostResponse, error) {
	db := p.db.Get()

	var post models.Post

	err := db.Preload("Author").First(&post, postID).Error
	if err != nil {
		if goerrors.Is(err, gorm.ErrRecordNotFound) {
			return PostResponse{}, errors.ErrNotFound
		}
		return PostResponse{}, err
	}
	return MapPostToResponse(post), nil
}

func (p *PostService) GetPosts(filter database.Filter) ([]PostResponse, error) {
	db := p.db.Get()

	var posts []models.Post

	query := db.Preload("Author")
	query = filter.Apply(query)

	err := query.Find(&posts).Error
	if err != nil {
		return nil, err
	}

	return MapPostsToResponse(posts), nil
}
