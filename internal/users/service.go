package users

import (
	"blog-api/internal/database"
	"blog-api/internal/errors"
	"blog-api/internal/logger"
	"blog-api/internal/models"
	"context"
	goerrors "errors"
	"log/slog"
)

type IUserService interface {
	GetUserByID(ctx context.Context, userID uint) (*UserResponse, error)
}

type userService struct {
	db     *database.DB
	logger *slog.Logger
}

func NewUserService(db *database.DB, logger *slog.Logger) IUserService {
	return &userService{
		db:     db,
		logger: logger,
	}
}

func (s *userService) GetUserByID(ctx context.Context, userID uint) (*UserResponse, error) {
	db := s.db.WithContext(ctx)
	log := logger.WithUserID(logger.FromCtx(ctx, s.logger), userID)

	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		if goerrors.Is(err, database.ErrRecordNotFound) {
			log.Warn("user not found")
			return nil, errors.ErrNotFound
		}

		log.Error("database query failed", logger.Err(err))
		return nil, err
	}
	return MapUserToResponse(user), nil
}
