package config

import "time"

type ServerConfig struct {
	Host            string        `validate:"required"`
	Port            string        `validate:"required"`
	ReadTimeout     time.Duration `validate:"required"`
	WriteTimeout    time.Duration `validate:"required"`
	ShutdownTimeout time.Duration `validate:"required"`
}
