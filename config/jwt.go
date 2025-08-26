package config

import "time"

type JwtConfig struct {
	AccessTTL  time.Duration `validate:"required"`
	RefreshTTL time.Duration `validate:"required"`
}
