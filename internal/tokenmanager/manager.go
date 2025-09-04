package tokenmanager

import "time"

type TokenManager interface {
	GenerateToken(userID string, tokenType string, ttl ...time.Duration) (string, time.Duration, error)
	ValidateToken(tokenStr string) (*Claims, error)
}
