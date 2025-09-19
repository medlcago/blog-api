package config

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	SecretKey string `validate:"required"`
	Env       string `validate:"required,oneof=local dev prod"`

	ServerConfig   ServerConfig   `validate:"required"`
	DatabaseConfig DatabaseConfig `validate:"required"`
	JwtConfig      JwtConfig      `validate:"required"`
	RedisConfig    RedisConfig    `validate:"required"`
	MinioConfig    MinioConfig    `validate:"required"`
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
		SecretKey:      v.GetString("SECRET_KEY"),
		Env:            v.GetString("APP_ENV"),
		ServerConfig:   loadServerConfig(v),
		DatabaseConfig: loadDatabaseConfig(v),
		JwtConfig:      loadJWTConfig(v),
		RedisConfig:    loadRedisConfig(v),
		MinioConfig:    loadMinioConfig(v),
	}

	if err := validateConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

func validateConfig(cfg *Config) error {
	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return fmt.Errorf("missing required attributes: %w", err)
	}
	return nil
}
