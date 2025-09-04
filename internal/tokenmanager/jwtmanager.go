package tokenmanager

import (
	"blog-api/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	AccessToken  = "access"
	RefreshToken = "refresh"

	DefaultTokenTTL = 10 * time.Minute
)

type JWTManager struct {
	secretKey  []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewJWTManager(secret string, cfg config.JwtConfig) TokenManager {
	return &JWTManager{
		secretKey:  []byte(secret),
		accessTTL:  cfg.AccessTTL,
		refreshTTL: cfg.RefreshTTL,
	}
}

func (m *JWTManager) GenerateToken(userID string, tokenType string, ttl ...time.Duration) (string, time.Duration, error) {
	tokenTTL := DefaultTokenTTL
	switch tokenType {
	case AccessToken:
		tokenTTL = m.accessTTL
	case RefreshToken:
		tokenTTL = m.refreshTTL
	default:
		if len(ttl) > 0 {
			tokenTTL = ttl[0]
		}
	}

	now := time.Now().UTC()
	claims := &Claims{
		TokenType: tokenType,
		UserID:    userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			Subject:   userID,
			ID:        uuid.NewString(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString(m.secretKey)
	if err != nil {
		return "", 0, err
	}
	return tokenStr, tokenTTL, nil
}

func (m *JWTManager) ValidateToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return m.secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}
	return claims, nil
}
