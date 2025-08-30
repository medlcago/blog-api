package auth

import (
	"blog-api/internal/database"
	"blog-api/internal/errors"
	"blog-api/internal/logger"
	"blog-api/internal/models"
	"blog-api/internal/storage"
	"blog-api/internal/tokenmanager"
	"blog-api/pkg/password"
	"context"
	goerrors "errors"
	"fmt"
	"log/slog"
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
	jwtService tokenmanager.JWTService
	db         *database.DB
	redis      *storage.RedisClient
	logger     *slog.Logger
}

func NewAuthService(jwtService tokenmanager.JWTService, db *database.DB, redis *storage.RedisClient, logger *slog.Logger) IAuthService {
	return &AuthService{
		jwtService: jwtService,
		db:         db,
		redis:      redis,
		logger:     logger,
	}
}

func (s *AuthService) Token(userID string) (*TokenResponse, error) {
	accessToken, err1 := s.jwtService.GenerateToken(userID, tokenmanager.AccessToken)
	refreshToken, err2 := s.jwtService.GenerateToken(userID, tokenmanager.RefreshToken)

	if err := goerrors.Join(err1, err2); err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:           accessToken,
		AccessTokenExpiresIn:  int(s.jwtService.AccessTTL().Seconds()),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresIn: int(s.jwtService.RefreshTTL().Seconds()),
	}, nil
}

func (s *AuthService) Register(ctx context.Context, input RegisterUserInput) (*TokenResponse, error) {
	db := s.db.Get().WithContext(ctx)
	log := logger.FromCtx(ctx, s.logger).With(slog.String("username", input.Username))

	log.Info("register attempt")

	normalizedUsername := strings.ToLower(strings.TrimSpace(input.Username))
	key := fmt.Sprintf("auth:register_attempt:%s", normalizedUsername)
	ttl := 24 * 7 * time.Hour

	set, err := s.redis.Client.SetNX(ctx, key, true, ttl).Result()
	if err != nil {
		log.Error("failed to set redis key", slog.String("key", key), logger.Err(err))
		return nil, err
	}
	if !set {
		log.Info("username already exists in redis")
		return nil, errors.ErrUsernameAlreadyExists
	}

	var user models.User
	err = db.Where("LOWER(username) = LOWER(?)", input.Username).First(&user).Error

	if err == nil {
		log.Info("username already exists in database")
		return nil, errors.ErrUsernameAlreadyExists
	}
	if !goerrors.Is(err, database.ErrRecordNotFound) {
		log.Error("database query failed", logger.Err(err))
		return nil, err
	}

	hashedPassword, err := password.HashPassword(input.Password)
	if err != nil {
		log.Error("password hashing failed", logger.Err(err))
		return nil, err
	}

	user = models.User{
		Username: input.Username,
		Password: hashedPassword,
	}

	if err = db.Create(&user).Error; err != nil {
		log.Error("user creation failed", logger.Err(err))
		return nil, err
	}

	userID := strconv.Itoa(int(user.ID))
	token, err := s.Token(userID)
	if err != nil {
		log.Error("failed to generate token", logger.Err(err))
		return nil, err
	}

	log.Info("user registered successfully")

	return token, nil

}

func (s *AuthService) Login(ctx context.Context, input LoginUserInput) (*TokenResponse, error) {
	db := s.db.Get().WithContext(ctx)
	log := logger.FromCtx(ctx, s.logger).With(slog.String("username", input.Username))

	log.Info("login attempt")

	var user models.User

	if err := db.Where("LOWER(username) = LOWER(?)", input.Username).First(&user).Error; err != nil {
		if !goerrors.Is(err, database.ErrRecordNotFound) {
			log.Error("database query failed", logger.Err(err))
			return nil, err
		}

		log.Info("user not found in database")
		return nil, errors.ErrInvalidCredentials
	}

	if !password.CheckPasswordHash(input.Password, user.Password) {
		log.Info("invalid password provided")
		return nil, errors.ErrInvalidCredentials
	}

	userID := strconv.Itoa(int(user.ID))
	token, err := s.Token(userID)
	if err != nil {
		log.Error("failed to generate token", logger.Err(err))
		return nil, err
	}

	log.Info("user successfully logged in")

	return token, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, input RefreshTokenInput) (*TokenResponse, error) {
	log := logger.FromCtx(ctx, s.logger)

	claims, err := s.jwtService.ValidateToken(input.RefreshToken)
	if err != nil {
		log.Info("invalid token provided", logger.Err(err))
		return nil, errors.ErrInvalidToken
	}

	if claims.TokenType != tokenmanager.RefreshToken {
		log.Info("wrong token type provided",
			slog.String("expected", tokenmanager.RefreshToken),
			slog.String("got", claims.TokenType),
		)
		return nil, errors.ErrInvalidToken
	}

	ttl := claims.GetRemainingDuration()
	if ttl <= 0 {
		log.Info("token has expired")
		return nil, errors.ErrInvalidToken
	}

	key := fmt.Sprintf("auth:revoked_token:%s", claims.ID)
	set, err := s.redis.Client.SetNX(ctx, key, "1", ttl).Result()
	if err != nil {
		log.Error("failed to set redis key", slog.String("key", key), logger.Err(err))
		return nil, err
	}

	if !set {
		log.Info("token has already been revoked")
		return nil, errors.ErrInvalidToken
	}

	log.Info("token revoked successfully", slog.String("token_id", claims.ID))

	token, err := s.Token(claims.UserID)
	if err != nil {
		log.Error("failed to generate new token", logger.Err(err))
		return nil, err
	}

	log.Info("token successfully refreshed")

	return token, err
}
