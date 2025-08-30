package posts

import (
	"blog-api/internal/database"
	"blog-api/internal/errors"
	"blog-api/internal/logger"
	"blog-api/internal/models"
	"context"
	goerrors "errors"
	"log/slog"

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
	db     *database.DB
	logger *slog.Logger
}

func NewPostService(db *database.DB, logger *slog.Logger) IPostService {
	return &PostService{
		db:     db,
		logger: logger,
	}
}

func (s *PostService) CreatePost(ctx context.Context, userID uint, input CreatePostInput) (*PostResponse, error) {
	db := s.db.Get().WithContext(ctx)
	log := logger.FromCtx(ctx, s.logger).With(slog.Any("user_id", userID))

	log.Info("creating new post")
	log.Debug("post input data", slog.Any("input", input))

	post := models.Post{
		Title:    input.Title,
		Content:  input.Content,
		AuthorID: userID,
	}

	entities := MapInputsToPostEntity(input.Entities)
	post.Entities = entities
	if err := ValidatePostEntities(post); err != nil {
		log.Warn("post entities validation failed", logger.Err(err))
		return nil, errors.New(400, err.Error())
	}

	if err := db.Create(&post).Error; err != nil {
		log.Error("failed to create post in database", logger.Err(err))
		return nil, err
	}

	if err := db.Preload("Author").Preload("Entities").First(&post, post.ID).Error; err != nil {
		log.Error("failed to load post with relations",
			logger.Err(err),
			slog.Any("post_id", post.ID),
		)
		return nil, err
	}

	log.Info("post created successfully", slog.Any("post_id", post.ID))

	return MapPostToResponse(post), nil
}

func (s *PostService) GetPost(ctx context.Context, postID uint) (*PostResponse, error) {
	db := s.db.Get().WithContext(ctx)
	log := logger.FromCtx(ctx, s.logger).With(slog.Any("post_id", postID))

	log.Info("fetching post")

	var post models.Post

	err := db.Preload("Author").Preload("Entities").First(&post, postID).Error
	if err != nil {
		if goerrors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("post not found")
			return nil, errors.ErrNotFound
		}

		log.Error("failed to fetch post from database", logger.Err(err))
		return nil, err
	}

	log.Info("post retrieved successfully")

	return MapPostToResponse(post), nil
}

func (s *PostService) GetPosts(ctx context.Context, params FilterParams) (*ListResponse, error) {
	db := s.db.Get().WithContext(ctx)
	log := logger.FromCtx(ctx, s.logger)

	log.Info("fetching posts list")
	log.Debug("filter parameters", slog.Any("params", params))

	var posts []models.Post

	query := db.Preload("Author").Preload("Entities").
		Scopes(
			OrderScope(params.OrderBy, params.Sort),
			PaginationScope(params.Limit, params.Offset),
		)

	err := query.Find(&posts).Error
	if err != nil {
		log.Error("failed to fetch posts from database", logger.Err(err))
		return nil, err
	}

	var total int64
	if err := db.Model(&posts).Count(&total).Error; err != nil {
		log.Error("failed to count posts", logger.Err(err))
		return nil, err
	}

	result := MapPostsToResponse(posts)
	listResponse := &ListResponse{
		Total:  total,
		Result: result,
	}

	log.Info("posts retrieved successfully", slog.Int64("total", total))

	return listResponse, nil
}

func (s *PostService) UpdatePost(ctx context.Context, userID uint, postID uint, input CreatePostInput) (*PostResponse, error) {
	db := s.db.Get().WithContext(ctx)
	log := logger.FromCtx(
		ctx,
		s.logger,
	).With(slog.Any("user_id", userID), slog.Any("post_id", postID))

	log.Info("updating post")
	log.Debug("update input data", slog.Any("input", input))

	err := db.Transaction(func(tx *gorm.DB) error {
		var post models.Post
		if err := tx.Preload("Author").First(&post, postID).Error; err != nil {
			if goerrors.Is(err, database.ErrRecordNotFound) {
				log.Warn("post not found")
				return errors.ErrNotFound
			}
			log.Error("failed to fetch post for update", logger.Err(err))
			return err
		}

		if post.AuthorID != userID {
			log.Warn("unauthorized update attempt",
				slog.Any("post_author_id", post.AuthorID),
				slog.Any("request_user_id", userID),
			)
			return errors.New(403, "Forbidden: You are not author of this post")
		}

		log.Debug("updating post fields")
		if err := tx.Model(&post).Updates(models.Post{
			Title:   input.Title,
			Content: input.Content,
		}).Error; err != nil {
			log.Error("failed to update post fields", logger.Err(err))
			return err
		}

		log.Debug("clearing existing entities")
		if err := tx.Unscoped().Model(&post).Association("Entities").Unscoped().Clear(); err != nil {
			log.Error("failed to clear post entities", logger.Err(err))
			return err
		}

		if len(input.Entities) > 0 {
			entities := MapInputsToPostEntity(input.Entities)
			post.Entities = entities

			if err := ValidatePostEntities(post); err != nil {
				log.Error("post entities validation failed during update", logger.Err(err))
				return errors.New(400, err.Error())
			}

			for i := range entities {
				entities[i].PostID = post.ID
			}

			log.Debug("creating new entities")
			if err := tx.Create(&entities).Error; err != nil {
				log.Error("failed to create new entities", logger.Err(err))
				return err
			}
		}

		return nil
	})

	if err != nil {
		log.Error("transaction failed during post update", logger.Err(err))
		return nil, err
	}

	var updatedPost models.Post
	if err = db.Preload("Author").Preload("Entities").First(&updatedPost, postID).Error; err != nil {
		log.Error("failed to load updated post", logger.Err(err))
		return nil, err
	}

	log.Info("post updated successfully")

	return MapPostToResponse(updatedPost), nil
}

func (s *PostService) DeletePost(ctx context.Context, userID uint, postID uint) error {
	db := s.db.Get().WithContext(ctx)
	log := logger.FromCtx(
		ctx,
		s.logger,
	).With(slog.Any("user_id", userID), slog.Any("post_id", postID))

	log.Info("deleting post")

	var post models.Post
	if err := db.First(&post, postID).Error; err != nil {
		if goerrors.Is(err, database.ErrRecordNotFound) {
			log.Warn("post not found")
			return errors.ErrNotFound
		}
		log.Error("failed to fetch post for deletion", logger.Err(err))
		return err
	}

	if post.AuthorID != userID {
		log.Warn("unauthorized deletion attempt",
			"post_author_id", post.AuthorID,
			"request_user_id", userID,
		)
		return errors.New(403, "Forbidden: You are not author of this post")
	}

	if err := db.Unscoped().Select("Entities").Delete(&post).Error; err != nil {
		log.Error("failed to delete post", logger.Err(err))
		return err
	}

	log.Info("post deleted successfully")

	return nil
}
