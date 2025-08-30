package users

import (
	"blog-api/internal/database"
	"blog-api/internal/logger"
	"blog-api/internal/models"
	"context"
	goerrors "errors"
	"log/slog"
)

type IUserService interface {
	GetUserByID(ctx context.Context, userID uint) (*UserResponse, error)
}

type UserService struct {
	db     *database.DB
	logger *slog.Logger
}

func NewUserService(db *database.DB, logger *slog.Logger) IUserService {
	return &UserService{
		db:     db,
		logger: logger,
	}
}

func (s *UserService) GetUserByID(ctx context.Context, userID uint) (*UserResponse, error) {
	db := s.db.Get().WithContext(ctx)
	log := logger.FromCtx(ctx, s.logger).With(slog.Any("user_id", userID))

	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		if goerrors.Is(err, database.ErrRecordNotFound) {
			log.Warn("user not found")
		}

		log.Error("database query failed", logger.Err(err))
		return nil, err
	}
	return MapUserToResponse(user), nil
}
