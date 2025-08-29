package tokenmanager

import "time"

type TokenService interface {
	GenerateToken(userID string, tokenType string) (string, error)
	ValidateToken(tokenStr string) (*Claims, error)
}

type JWTService interface {
	TokenService
	AccessTTL() time.Duration
	RefreshTTL() time.Duration
}
