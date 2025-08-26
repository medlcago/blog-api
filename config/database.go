package config

import "time"

type DatabaseConfig struct {
	Host     string `validate:"required"`
	User     string `validate:"required"`
	Password string `validate:"required"`
	Name     string `validate:"required"`
	Port     string `validate:"required"`
	SSLMode  string `validate:"required"`
	TimeZone string `validate:"required"`

	MaxIdleConns    int           `validate:"required"`
	MaxOpenConns    int           `validate:"required"`
	ConnMaxLifetime time.Duration `validate:"required"`
}
