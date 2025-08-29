package users

import (
	"blog-api/internal/database"
	"blog-api/internal/models"
	"context"
)

type IUserService interface {
	GetUserByID(ctx context.Context, userID uint) (*UserResponse, error)
}

type UserService struct {
	db *database.DB
}

func NewUserService(db *database.DB) IUserService {
	return &UserService{
		db: db,
	}
}

func (s *UserService) GetUserByID(ctx context.Context, userID uint) (*UserResponse, error) {
	db := s.db.Get().WithContext(ctx)

	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return MapUserToResponse(user), nil
}
