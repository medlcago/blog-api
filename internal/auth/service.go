package auth

import (
	"blog-api/internal/database"
	"blog-api/internal/errors"
	"blog-api/internal/logger"
	"blog-api/internal/models"
	"blog-api/internal/storage"
	"blog-api/internal/tokenmanager"
	"blog-api/pkg/password"
	"bytes"
	"context"
	"encoding/base64"
	goerrors "errors"
	"fmt"
	"image/png"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/pquerna/otp/totp"
)

const (
	registerAttemptKeyPrefix = "register_attempt"
	twoFAAuthKeyPrefix       = "2fa_auth"
	refreshTokenKeyPrefix    = "refresh_token"
)

func (s *authService) getRegisterAttemptKey(username string) string {
	normalizedUsername := strings.ToLower(strings.TrimSpace(username))
	return fmt.Sprintf("%s:%s", registerAttemptKeyPrefix, normalizedUsername)
}

func (s *authService) get2FAAuthKey(userID uint) string {
	return fmt.Sprintf("%s:%d", twoFAAuthKeyPrefix, userID)
}

func (s *authService) getRefreshTokenKey(tokenID string) string {
	return fmt.Sprintf("%s:%s", refreshTokenKeyPrefix, tokenID)
}

type IAuthService interface {
	Register(ctx context.Context, input RegisterUserInput) (*TokenResponse, error)
	Login(ctx context.Context, input LoginUserInput) (*LoginResponse, error)
	RefreshToken(ctx context.Context, input RefreshTokenInput) (*TokenResponse, error)
	ChangePassword(ctx context.Context, userID uint, input ChangePasswordInput) error

	Login2FA(ctx context.Context, input Login2FAInput) (*TokenResponse, error)
	Enable2FA(ctx context.Context, userID uint) (*TwoFASetupResponse, error)
	Verify2FA(ctx context.Context, userID uint, input Verify2FAInput) error
	Disable2FA(ctx context.Context, userID uint) error
}

type authService struct {
	tokenService tokenmanager.TokenManager
	db           *database.DB
	redis        *storage.RedisClient
	logger       *slog.Logger
}

func NewAuthService(tokenService tokenmanager.TokenManager, db *database.DB, redis *storage.RedisClient, logger *slog.Logger) IAuthService {
	return &authService{
		tokenService: tokenService,
		db:           db,
		redis:        redis,
		logger:       logger,
	}
}

func (s *authService) Token(userID string) (*TokenResponse, error) {
	accessToken, accessTokenTTL, err1 := s.tokenService.GenerateToken(userID, tokenmanager.AccessToken)
	refreshToken, refreshTokenTTL, err2 := s.tokenService.GenerateToken(userID, tokenmanager.RefreshToken)

	if err := goerrors.Join(err1, err2); err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:           accessToken,
		AccessTokenExpiresIn:  int(accessTokenTTL.Seconds()),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresIn: int(refreshTokenTTL.Seconds()),
	}, nil
}

func (s *authService) Register(ctx context.Context, input RegisterUserInput) (*TokenResponse, error) {
	db := s.db.WithContext(ctx)
	log := logger.WithUsername(logger.FromCtx(ctx, s.logger), input.Username)

	log.Info("register attempt")

	key := s.getRegisterAttemptKey(input.Username)
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

func (s *authService) Login(ctx context.Context, input LoginUserInput) (*LoginResponse, error) {
	db := s.db.WithContext(ctx)
	log := logger.WithUsername(logger.FromCtx(ctx, s.logger), input.Username)

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

	if user.TwoFAEnabled {
		log.Info("2FA required for user")

		key := s.get2FAAuthKey(user.ID)
		err := s.redis.Client.SetNX(ctx, key, "1", 5*time.Minute).Err()
		if err != nil {
			log.Error("failed to set 2FA auth stage", logger.Err(err))
			return nil, err
		}

		return &LoginResponse{
			Requires2FA: true,
			Message:     "2FA code required",
		}, nil
	}

	userID := strconv.Itoa(int(user.ID))
	token, err := s.Token(userID)
	if err != nil {
		log.Error("failed to generate token", logger.Err(err))
		return nil, err
	}

	log.Info("user successfully logged in (without 2FA)")

	return &LoginResponse{
		Token:       token,
		Requires2FA: false,
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, input RefreshTokenInput) (*TokenResponse, error) {
	log := logger.FromCtx(ctx, s.logger)

	claims, err := s.tokenService.ValidateToken(input.RefreshToken)
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

	key := s.getRefreshTokenKey(claims.ID)
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

func (s *authService) ChangePassword(ctx context.Context, userID uint, input ChangePasswordInput) error {
	db := s.db.WithContext(ctx)
	log := logger.WithUserID(logger.FromCtx(ctx, s.logger), userID)

	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		if goerrors.Is(err, database.ErrRecordNotFound) {
			log.Warn("user not found")
			return errors.ErrUnauthorized
		}
		log.Error("failed to get user", logger.Err(err))
		return err
	}

	if !password.CheckPasswordHash(input.OldPassword, user.Password) {
		return errors.ErrIncorrectOldPassword
	}

	if input.OldPassword == input.NewPassword {
		log.Warn("new password matches old password")
		return errors.ErrNewPasswordSameAsOld
	}

	hashedPassword, err := password.HashPassword(input.NewPassword)
	if err != nil {
		log.Error("failed to hash password", logger.Err(err))
		return err
	}

	if err := db.Model(&user).Update("password", hashedPassword).Error; err != nil {
		log.Error("failed to change password", logger.Err(err))
		return err
	}

	return nil
}

func (s *authService) Login2FA(ctx context.Context, input Login2FAInput) (*TokenResponse, error) {
	db := s.db.WithContext(ctx)
	log := logger.WithUsername(logger.FromCtx(ctx, s.logger), input.Username)

	var user models.User
	if err := db.First(&user, "LOWER(username) = LOWER(?)", input.Username).Error; err != nil {
		if goerrors.Is(err, database.ErrRecordNotFound) {
			log.Warn("user not found")
			return nil, errors.ErrUnauthorized
		}
		log.Error("failed to get user", logger.Err(err))
		return nil, err
	}

	key := s.get2FAAuthKey(user.ID)
	exists, err := s.redis.Client.Exists(ctx, key).Result()
	if err != nil {
		log.Error("failed to check 2FA auth key key in redis", logger.Err(err))
		return nil, err
	}
	if exists == 0 {
		log.Warn("2FA auth key not found in redis, login flow not initiated")
		return nil, errors.ErrTwoFAFlowNotInitiated
	}

	if !user.TwoFAEnabled {
		log.Info("2FA not enabled for user")
		return nil, errors.ErrTwoFANotEnabled
	}

	if !totp.Validate(input.Code, user.TwoFASecret.String) {
		log.Warn("invalid 2FA code")
		return nil, errors.ErrInvalid2FACode
	}

	if err := s.redis.Client.Del(ctx, key).Err(); err != nil {
		log.Error("failed to delete 2FA auth key", logger.Err(err))
		return nil, err
	}

	userID := strconv.Itoa(int(user.ID))
	token, err := s.Token(userID)
	if err != nil {
		log.Error("failed to generate token", logger.Err(err))
		return nil, err
	}

	log.Info("2FA verification successful, user logged in")
	return token, err
}

func (s *authService) Enable2FA(ctx context.Context, userID uint) (*TwoFASetupResponse, error) {
	db := s.db.WithContext(ctx)
	log := logger.WithUserID(logger.FromCtx(ctx, s.logger), userID)

	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		if goerrors.Is(err, database.ErrRecordNotFound) {
			log.Warn("user not found")
			return nil, errors.ErrUnauthorized
		}
		log.Error("failed to get user", logger.Err(err))
		return nil, err
	}

	if user.TwoFAEnabled {
		log.Info("2FA is already enabled")
		return nil, errors.ErrTwoFAAlreadyEnabled
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "BlogAPI",
		AccountName: strings.ToLower(user.Username),
	})

	if err != nil {
		log.Error("failed to generate TOTP key", logger.Err(err))
		return nil, err
	}

	if err = db.Model(&user).Updates(map[string]any{
		"two_fa_secret": key.Secret(),
	}).Error; err != nil {
		log.Error("failed to set two_fa_secret", logger.Err(err))
		return nil, err
	}

	qrCode, err := key.Image(200, 200)
	if err != nil {
		log.Error("failed to generate QRCode", logger.Err(err))
		return nil, err
	}

	var buffer bytes.Buffer
	if err := png.Encode(&buffer, qrCode); err != nil {
		log.Error("failed to encode QRCode", logger.Err(err))
		return nil, err
	}

	return &TwoFASetupResponse{
		QRCode:  "data:image/png;base64," + base64.StdEncoding.EncodeToString(buffer.Bytes()),
		Message: "Scan QR code with authenticator app and verify with a code",
	}, nil
}

func (s *authService) Verify2FA(ctx context.Context, userID uint, input Verify2FAInput) error {
	db := s.db.WithContext(ctx)
	log := logger.WithUserID(logger.FromCtx(ctx, s.logger), userID)

	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		if goerrors.Is(err, database.ErrRecordNotFound) {
			log.Warn("user not found")
			return errors.ErrUnauthorized
		}
		log.Error("failed to get user", logger.Err(err))
		return err
	}

	if user.TwoFAEnabled {
		log.Info("2FA is already enabled")
		return errors.ErrTwoFAAlreadyEnabled
	}

	if !totp.Validate(input.Code, user.TwoFASecret.String) {
		return errors.ErrInvalid2FACode
	}

	if err := db.Model(&user).Updates(map[string]any{
		"two_fa_enabled": true,
	}).Error; err != nil {
		log.Error("failed to enable 2FA", logger.Err(err))
		return err
	}
	return nil
}

func (s *authService) Disable2FA(ctx context.Context, userID uint) error {
	db := s.db.WithContext(ctx)
	log := logger.WithUserID(logger.FromCtx(ctx, s.logger), userID)

	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		if goerrors.Is(err, database.ErrRecordNotFound) {
			log.Warn("user not found")
			return errors.ErrUnauthorized
		}
		log.Error("failed to get user", logger.Err(err))
		return err
	}

	if !user.TwoFAEnabled {
		log.Info("2FA is not enabled")
		return errors.ErrTwoFANotEnabled
	}

	if err := db.Model(&user).Updates(map[string]any{
		"two_fa_enabled": false,
		"two_fa_secret":  nil,
	}).Error; err != nil {
		log.Error("failed to disable 2FA", logger.Err(err))
		return err
	}

	return nil
}
