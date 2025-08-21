package config

import "time"

type JwtConfig struct {
	JwtAccessTTL  time.Duration `validate:"required"`
	JwtRefreshTTL time.Duration `validate:"required"`
}
