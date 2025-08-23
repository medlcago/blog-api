package posts

import (
	"blog-api/internal/database"
	"blog-api/internal/models"
	"blog-api/pkg/errors"
	goerrors "errors"

	"gorm.io/gorm"
)

type IPostService interface {
	CreatePost(userID uint, input CreatePostInput) (*PostResponse, error)
	GetPost(postID uint) (*PostResponse, error)
	GetPosts(filter database.Filter) ([]*PostResponse, error)
	UpdatePost(userID uint, postID uint, input CreatePostInput) (*PostResponse, error)
	DeletePost(userID uint, postID uint) error
}

type PostService struct {
	db *database.DB
}

func NewPostService(db *database.DB) IPostService {
	return &PostService{
		db: db,
	}
}

func (p *PostService) CreatePost(userID uint, input CreatePostInput) (*PostResponse, error) {
	db := p.db.Get()

	post := models.Post{
		Title:    input.Title,
		Content:  input.Content,
		AuthorID: userID,
	}

	entities := MapInputsToPostEntity(input.Entities)
	post.Entities = entities
	if err := ValidatePostEntities(post); err != nil {
		return nil, errors.New(400, err.Error())
	}

	if err := db.Create(&post).Error; err != nil {
		return nil, err
	}

	if err := db.Preload("Author").Preload("Entities").First(&post, post.ID).Error; err != nil {
		return nil, err
	}

	return MapPostToResponse(post), nil
}

func (p *PostService) GetPost(postID uint) (*PostResponse, error) {
	db := p.db.Get()

	var post models.Post

	err := db.Preload("Author").Preload("Entities").First(&post, postID).Error
	if err != nil {
		if goerrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.ErrNotFound
		}
		return nil, err
	}

	return MapPostToResponse(post), nil
}

func (p *PostService) GetPosts(filter database.Filter) ([]*PostResponse, error) {
	db := p.db.Get()

	var posts []models.Post

	query := db.Preload("Author").Preload("Entities")
	query = filter.Apply(query)

	err := query.Find(&posts).Error
	if err != nil {
		return nil, err
	}

	return MapPostsToResponse(posts), nil
}

func (p *PostService) UpdatePost(userID uint, postID uint, input CreatePostInput) (*PostResponse, error) {
	db := p.db.Get()

	err := db.Transaction(func(tx *gorm.DB) error {
		var post models.Post
		if err := tx.Preload("Author").First(&post, postID).Error; err != nil {
			if goerrors.Is(err, gorm.ErrRecordNotFound) {
				return errors.ErrNotFound
			}
			return err
		}

		if post.AuthorID != userID {
			return errors.New(403, "Forbidden: You are not author of this post")
		}

		if err := tx.Model(&post).Updates(models.Post{
			Title:   input.Title,
			Content: input.Content,
		}).Error; err != nil {
			return err
		}

		if err := tx.Unscoped().Model(&post).Association("Entities").Unscoped().Clear(); err != nil {
			return err
		}

		if len(input.Entities) > 0 {
			entities := MapInputsToPostEntity(input.Entities)
			post.Entities = entities
			if err := ValidatePostEntities(post); err != nil {
				return errors.New(400, err.Error())
			}

			for i := range entities {
				entities[i].PostID = post.ID
			}

			if err := tx.Create(&entities).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	var updatedPost models.Post
	if err = db.Preload("Author").Preload("Entities").First(&updatedPost, postID).Error; err != nil {
		return nil, err
	}

	return MapPostToResponse(updatedPost), nil
}

func (p *PostService) DeletePost(userID uint, postID uint) error {
	db := p.db.Get()

	var post models.Post
	if err := db.First(&post, postID).Error; err != nil {
		if goerrors.Is(err, gorm.ErrRecordNotFound) {
			return errors.ErrNotFound
		}
		return err
	}

	if post.AuthorID != userID {
		return errors.New(403, "Forbidden: You are not author of this post")
	}

	if err := db.Unscoped().Select("Entities").Delete(&post).Error; err != nil {
		return err
	}

	return nil
}
