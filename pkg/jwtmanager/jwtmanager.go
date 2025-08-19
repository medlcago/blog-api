package jwtmanager

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	AccessToken  = "access"
	RefreshToken = "refresh"
)

type Claims struct {
	TokenType string `json:"type"`
	UserID    string `json:"uid"`
	jwt.RegisteredClaims
}

type JWTManager struct {
	secretKey  []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewJWTManager(secret string, accessTTL, refreshTTL time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:  []byte(secret),
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

func (m *JWTManager) AccessTTL() time.Duration {
	return m.accessTTL
}

func (m *JWTManager) RefreshTTL() time.Duration {
	return m.refreshTTL
}

func (m *JWTManager) GenerateToken(userID string, tokenType string) (string, error) {
	var ttl time.Duration
	switch tokenType {
	case AccessToken:
		ttl = m.AccessTTL()
	case RefreshToken:
		ttl = m.RefreshTTL()
	default:
		return "", errors.New("unknown token type")
	}

	now := time.Now().UTC()
	claims := &Claims{
		TokenType: tokenType,
		UserID:    userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(now),
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secretKey)
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
		return nil, errors.New("invalid token claims")
	}
	return claims, nil
}
