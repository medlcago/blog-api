package config

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	SecretKey string `validate:"required"`

	ServerConfig   `validate:"required"`
	DatabaseConfig `validate:"required"`
	JwtConfig      `validate:"required"`
}

func MustGet() *Config {
	cfg, err := Get()
	if err != nil {
		panic("failed to get config: " + err.Error())
	}
	return cfg
}

func Get() (*Config, error) {
	v := viper.New()
	v.SetConfigFile(".env")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	v.AutomaticEnv()

	setDefaults(v)

	config := &Config{
		SecretKey: v.GetString("SECRET_KEY"),
		ServerConfig: ServerConfig{
			ServerHost: v.GetString("SERVER_HOST"),
			ServerPort: v.GetString("SERVER_PORT"),
		},
		DatabaseConfig: DatabaseConfig{
			DbHost:     v.GetString("DB_HOST"),
			DbUser:     v.GetString("DB_USER"),
			DbPassword: v.GetString("DB_PASSWORD"),
			DbName:     v.GetString("DB_NAME"),
			DbPort:     v.GetString("DB_PORT"),

			SSLMode:  v.GetString("SSL_MODE"),
			TimeZone: v.GetString("TIME_ZONE"),

			MaxIdleConns:    v.GetInt("MAX_IDLE_CONNS"),
			MaxOpenConns:    v.GetInt("MAX_OPEN_CONNS"),
			ConnMaxLifetime: v.GetDuration("CONN_MAX_LIFETIME"),
		},
		JwtConfig: JwtConfig{
			JwtAccessTTL:  v.GetDuration("JWT_ACCESS_TTL"),
			JwtRefreshTTL: v.GetDuration("JWT_REFRESH_TTL"),
		},
	}

	if err := validateConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("SERVER_HOST", "0.0.0.0")
	v.SetDefault("SERVER_PORT", "3000")

	v.SetDefault("SSL_MODE", "disable")
	v.SetDefault("TIME_ZONE", "UTC")

	v.SetDefault("MAX_IDLE_CONNS", 10)
	v.SetDefault("MAX_OPEN_CONNS", 100)
	v.SetDefault("CONN_MAX_LIFETIME", time.Hour)

	v.SetDefault("JWT_ACCESS_TTL", 30*time.Minute)
	v.SetDefault("JWT_REFRESH_TTL", 24*30*time.Hour)
}

func validateConfig(cfg *Config) error {
	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return fmt.Errorf("missing required attributes: %w", err)
	}
	return nil
}
