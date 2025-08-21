package config

import "time"

type DatabaseConfig struct {
	DbHost     string `validate:"required"`
	DbUser     string `validate:"required"`
	DbPassword string `validate:"required"`
	DbName     string `validate:"required"`
	DbPort     string `validate:"required"`

	MaxIdleConns    int           `validate:"required"`
	MaxOpenConns    int           `validate:"required"`
	ConnMaxLifetime time.Duration `validate:"required"`
}
