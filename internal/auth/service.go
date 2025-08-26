package auth

import (
	"blog-api/internal/database"
	"blog-api/internal/jwtmanager"
	"blog-api/internal/models"
	"blog-api/pkg/errors"
	"blog-api/pkg/password"
	"blog-api/pkg/storage"
	"context"
	goerrors "errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

type IAuthService interface {
	Register(ctx context.Context, input RegisterUserInput) (*TokenResponse, error)
	Login(ctx context.Context, input LoginUserInput) (*TokenResponse, error)
}

type AuthService struct {
	jwtManager *jwtmanager.JWTManager
	db         *database.DB
	store      storage.Storage
	appLogger  *log.Logger
}

func NewAuthService(jwtManager *jwtmanager.JWTManager, db *database.DB, store storage.Storage, appLogger *log.Logger) *AuthService {
	return &AuthService{
		jwtManager: jwtManager,
		db:         db,
		store:      store,
		appLogger:  appLogger,
	}
}

func (a *AuthService) Token(userID string) (*TokenResponse, error) {
	accessToken, err1 := a.jwtManager.GenerateToken(userID, jwtmanager.AccessToken)
	refreshToken, err2 := a.jwtManager.GenerateToken(userID, jwtmanager.RefreshToken)

	if err := goerrors.Join(err1, err2); err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:           accessToken,
		AccessTokenExpiresIn:  int(a.jwtManager.AccessTTL().Seconds()),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresIn: int(a.jwtManager.RefreshTTL().Seconds()),
	}, nil
}

func (a *AuthService) Register(ctx context.Context, input RegisterUserInput) (*TokenResponse, error) {
	db := a.db.Get().WithContext(ctx)

	normalizedUsername := strings.ToLower(strings.TrimSpace(input.Username))
	key := fmt.Sprintf("auth:register_attempt:%s", normalizedUsername)
	ttl := 24 * 7 * time.Hour

	exists, err := a.store.Exists(ctx, key)
	if err != nil {
		a.appLogger.Printf("storage check failed: %v", err)
	} else if exists {
		return nil, errors.ErrUsernameAlreadyExists
	}

	var user models.User
	err = db.Where("LOWER(username) = LOWER(?)", input.Username).First(&user).Error

	if err == nil {
		if err := a.store.Set(ctx, key, true, ttl); err != nil {
			a.appLogger.Printf("failed to storage username existence: %v", err)
		}
		return nil, errors.ErrUsernameAlreadyExists
	}

	if !goerrors.Is(err, database.ErrRecordNotFound) {
		return nil, err
	}

	hashedPassword, err := password.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	user = models.User{
		Username: input.Username,
		Password: hashedPassword,
	}
	if err = db.Create(&user).Error; err != nil {
		return nil, err
	}

	userID := strconv.Itoa(int(user.ID))
	token, err := a.Token(userID)
	if err != nil {
		return nil, err
	}

	if err := a.store.Set(ctx, key, true, ttl); err != nil {
		a.appLogger.Printf("failed to storage registration: %v", err)
	}

	return token, nil

}

func (a *AuthService) Login(ctx context.Context, input LoginUserInput) (*TokenResponse, error) {
	db := a.db.Get().WithContext(ctx)

	var user models.User

	if err := db.Where("LOWER(username) = LOWER(?)", input.Username).First(&user).Error; err != nil {
		return nil, errors.ErrInvalidCredentials
	}

	if !password.CheckPasswordHash(input.Password, user.Password) {
		return nil, errors.ErrInvalidCredentials
	}

	userID := strconv.Itoa(int(user.ID))
	return a.Token(userID)
}
