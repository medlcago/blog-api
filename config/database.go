package config

import (
	"time"

	"github.com/spf13/viper"
)

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

func loadDatabaseConfig(v *viper.Viper) DatabaseConfig {
	return DatabaseConfig{
		Host:     v.GetString("DB_HOST"),
		User:     v.GetString("DB_USER"),
		Password: v.GetString("DB_PASSWORD"),
		Name:     v.GetString("DB_NAME"),
		Port:     v.GetString("DB_PORT"),

		SSLMode:  v.GetString("DB_SSL_MODE"),
		TimeZone: v.GetString("DB_TIME_ZONE"),

		MaxIdleConns:    v.GetInt("DB_MAX_IDLE_CONNS"),
		MaxOpenConns:    v.GetInt("DB_MAX_OPEN_CONNS"),
		ConnMaxLifetime: v.GetDuration("DB_CONN_MAX_LIFETIME"),
	}
}
