package posts

import (
	"blog-api/internal/database"
	"blog-api/internal/models"
	"blog-api/pkg/errors"
	"context"
	goerrors "errors"

	"gorm.io/gorm"
)

type IPostService interface {
	CreatePost(ctx context.Context, userID uint, input CreatePostInput) (*PostResponse, error)
	GetPost(ctx context.Context, postID uint) (*PostResponse, error)
	GetPosts(ctx context.Context, params FilterParams) (*ListResponse, error)
	UpdatePost(ctx context.Context, userID uint, postID uint, input CreatePostInput) (*PostResponse, error)
	DeletePost(ctx context.Context, userID uint, postID uint) error
}

type PostService struct {
	db *database.DB
}

func NewPostService(db *database.DB) IPostService {
	return &PostService{
		db: db,
	}
}

func (p *PostService) CreatePost(ctx context.Context, userID uint, input CreatePostInput) (*PostResponse, error) {
	db := p.db.Get().WithContext(ctx)

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

func (p *PostService) GetPost(ctx context.Context, postID uint) (*PostResponse, error) {
	db := p.db.Get().WithContext(ctx)

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

func (p *PostService) GetPosts(ctx context.Context, params FilterParams) (*ListResponse, error) {
	db := p.db.Get().WithContext(ctx)

	var posts []models.Post

	query := db.Preload("Author").Preload("Entities").
		Scopes(
			OrderScope(params.OrderBy, params.Sort),
			PaginationScope(params.Limit, params.Offset),
		)

	err := query.Find(&posts).Error
	if err != nil {
		return nil, err
	}

	var total int64
	if err := db.Model(&posts).Count(&total).Error; err != nil {
		return nil, err
	}

	result := MapPostsToResponse(posts)
	listResponse := &ListResponse{
		Total:  total,
		Result: result,
	}

	return listResponse, nil
}

func (p *PostService) UpdatePost(ctx context.Context, userID uint, postID uint, input CreatePostInput) (*PostResponse, error) {
	db := p.db.Get().WithContext(ctx)

	err := db.Transaction(func(tx *gorm.DB) error {
		var post models.Post
		if err := tx.Preload("Author").First(&post, postID).Error; err != nil {
			if goerrors.Is(err, database.ErrRecordNotFound) {
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

func (p *PostService) DeletePost(ctx context.Context, userID uint, postID uint) error {
	db := p.db.Get().WithContext(ctx)

	var post models.Post
	if err := db.First(&post, postID).Error; err != nil {
		if goerrors.Is(err, database.ErrRecordNotFound) {
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
