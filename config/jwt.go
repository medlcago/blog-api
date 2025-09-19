package config

import (
	"time"

	"github.com/spf13/viper"
)

type JwtConfig struct {
	AccessTTL  time.Duration `validate:"required"`
	RefreshTTL time.Duration `validate:"required"`
}

func loadJWTConfig(v *viper.Viper) JwtConfig {
	return JwtConfig{
		AccessTTL:  v.GetDuration("JWT_ACCESS_TTL"),
		RefreshTTL: v.GetDuration("JWT_REFRESH_TTL"),
	}
}
