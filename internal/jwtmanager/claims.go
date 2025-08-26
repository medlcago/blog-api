package jwtmanager

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	TokenType string `json:"type"`
	UserID    string `json:"uid"`
	jwt.RegisteredClaims
}

func (c *Claims) GetDuration() time.Duration {
	return c.ExpiresAt.Time.Sub(c.IssuedAt.Time)
}

func (c *Claims) GetRemainingDuration() time.Duration {
	if c.ExpiresAt == nil {
		return 0
	}

	now := time.Now().UTC()
	remaining := c.ExpiresAt.Time.Sub(now)

	if remaining < 0 {
		return 0
	}

	return remaining
}
