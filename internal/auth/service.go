package auth

import (
	"blog-api/internal/database"
	"blog-api/internal/models"
	"blog-api/pkg/errors"
	"blog-api/pkg/jwtmanager"
	"blog-api/pkg/password"
	goerrors "errors"
	"strconv"
)

type IAuthService interface {
	Register(input RegisterUserInput) (TokenResponse, error)
	Login(input LoginUserInput) (TokenResponse, error)
}

type AuthService struct {
	jwtManager *jwtmanager.JWTManager
}

func NewAuthService(jwtManager *jwtmanager.JWTManager) *AuthService {
	return &AuthService{
		jwtManager: jwtManager,
	}
}

func (a *AuthService) Token(userID string) (TokenResponse, error) {
	accessToken, err1 := a.jwtManager.GenerateToken(userID, jwtmanager.AccessToken)
	refreshToken, err2 := a.jwtManager.GenerateToken(userID, jwtmanager.RefreshToken)

	if err := goerrors.Join(err1, err2); err != nil {
		return TokenResponse{}, err
	}

	return TokenResponse{
		AccessToken:           accessToken,
		AccessTokenExpiresIn:  int64(a.jwtManager.AccessTTL().Seconds()),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresIn: int64(a.jwtManager.RefreshTTL().Seconds()),
	}, nil
}

func (a *AuthService) Register(input RegisterUserInput) (TokenResponse, error) {
	instance := database.GetDb()

	var user models.User
	err := instance.Where("LOWER(username) = LOWER(?)", input.Username).First(&user).Error

	if err == nil {
		return TokenResponse{}, errors.ErrUsernameAlreadyExists
	}

	if !goerrors.Is(err, database.ErrRecordNotFound) {
		return TokenResponse{}, err
	}

	hashedPassword, err := password.HashPassword(input.Password)
	if err != nil {
		return TokenResponse{}, err
	}

	user = models.User{
		Username: input.Username,
		Password: hashedPassword,
	}
	if err = instance.Create(&user).Error; err != nil {
		return TokenResponse{}, err
	}

	userID := strconv.Itoa(int(user.ID))
	return a.Token(userID)

}

func (a *AuthService) Login(input LoginUserInput) (TokenResponse, error) {
	instance := database.GetDb()

	var user models.User

	if err := instance.Where("LOWER(username) = LOWER(?)", input.Username).First(&user).Error; err != nil {
		return TokenResponse{}, errors.ErrInvalidCredentials
	}

	if !password.CheckPasswordHash(input.Password, user.Password) {
		return TokenResponse{}, errors.ErrInvalidCredentials
	}

	userID := strconv.Itoa(int(user.ID))
	return a.Token(userID)
}
