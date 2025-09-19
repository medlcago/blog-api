package config

import (
	"time"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Host            string        `validate:"required"`
	Port            string        `validate:"required"`
	BodyLimit       int           `validate:"required"`
	ReadTimeout     time.Duration `validate:"required"`
	WriteTimeout    time.Duration `validate:"required"`
	ShutdownTimeout time.Duration `validate:"required"`
}

func loadServerConfig(v *viper.Viper) ServerConfig {
	return ServerConfig{
		Host:            v.GetString("SERVER_HOST"),
		Port:            v.GetString("SERVER_PORT"),
		BodyLimit:       v.GetInt("SERVER_BODY_LIMIT"),
		ReadTimeout:     v.GetDuration("SERVER_READ_TIMEOUT"),
		WriteTimeout:    v.GetDuration("SERVER_WRITE_TIMEOUT"),
		ShutdownTimeout: v.GetDuration("SERVER_SHUTDOWN_TIMEOUT"),
	}
}
