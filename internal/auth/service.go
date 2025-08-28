package auth

import (
	"blog-api/internal/database"
	"blog-api/internal/jwtmanager"
	"blog-api/internal/models"
	"blog-api/internal/storage"
	"blog-api/pkg/errors"
	"blog-api/pkg/password"
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
	RefreshToken(ctx context.Context, input RefreshTokenInput) (*TokenResponse, error)
}

type AuthService struct {
	jwtManager *jwtmanager.JWTManager
	db         *database.DB
	redis      *storage.RedisClient
	appLogger  *log.Logger
}

func NewAuthService(jwtManager *jwtmanager.JWTManager, db *database.DB, redis *storage.RedisClient, appLogger *log.Logger) IAuthService {
	return &AuthService{
		jwtManager: jwtManager,
		db:         db,
		redis:      redis,
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

	exists, err := a.redis.Client.Exists(ctx, key).Result()
	if err != nil {
		a.appLogger.Printf("redis check failed for key %s: %v", key, err)
	} else if exists == 1 {
		return nil, errors.ErrUsernameAlreadyExists
	}

	var user models.User
	err = db.Where("LOWER(username) = LOWER(?)", input.Username).First(&user).Error

	if err == nil {
		if err := a.redis.Client.Set(ctx, key, true, ttl).Err(); err != nil {
			a.appLogger.Printf("failed to set redis key %s: %v", key, err)
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

	if err := a.redis.Client.Set(ctx, key, true, ttl).Err(); err != nil {
		a.appLogger.Printf("failed to set redis key %s: %v", key, err)
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

func (a *AuthService) RefreshToken(ctx context.Context, input RefreshTokenInput) (*TokenResponse, error) {
	claims, err := a.jwtManager.ValidateToken(input.RefreshToken)
	if err != nil {
		return nil, errors.ErrInvalidToken
	}

	if claims.TokenType != jwtmanager.RefreshToken {
		return nil, errors.ErrInvalidToken
	}

	ttl := claims.GetRemainingDuration()
	if ttl <= 0 {
		return nil, errors.ErrInvalidToken
	}

	key := fmt.Sprintf("auth:revoked_token:%s", claims.ID)
	exists, err := a.redis.Client.Exists(ctx, key).Result()
	if exists == 1 && err == nil {
		return nil, errors.ErrInvalidToken
	}

	if err := a.redis.Client.Set(ctx, key, "1", ttl).Err(); err != nil {
		return nil, err
	}

	return a.Token(claims.UserID)
}
