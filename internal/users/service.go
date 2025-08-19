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
}

func NewUserService(jwtManager *jwtmanager.JWTManager) IUserService {
	return &UserService{
		jwtManager: jwtManager,
	}
}

func (s *UserService) GetUserByID(userID uint) (models.User, error) {
	instance := database.GetDb()

	var user models.User
	if err := instance.First(&user, userID).Error; err != nil {
		return models.User{}, err
	}
	return user, nil
}
