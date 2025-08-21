package users

import (
	"blog-api/internal/database"
	"blog-api/internal/models"
	"blog-api/pkg/jwtmanager"
)

type IUserService interface {
	GetUserByID(userID uint) (models.User, error)
}

type UserService struct {
	jwtManager *jwtmanager.JWTManager
	db         *database.DB
}

func NewUserService(jwtManager *jwtmanager.JWTManager, db *database.DB) IUserService {
	return &UserService{
		jwtManager: jwtManager,
		db:         db,
	}
}

func (s *UserService) GetUserByID(userID uint) (models.User, error) {
	db := s.db.Get()

	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		return models.User{}, err
	}
	return user, nil
}
