package config

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	SecretKey string `validate:"required"`

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
		SecretKey: v.GetString("SECRET_KEY"),
		ServerConfig: ServerConfig{
			Host: v.GetString("SERVER_HOST"),
			Port: v.GetString("SERVER_PORT"),
		},
		DatabaseConfig: DatabaseConfig{
			Host:     v.GetString("DB_HOST"),
			User:     v.GetString("DB_USER"),
			Password: v.GetString("DB_PASSWORD"),
			Name:     v.GetString("DB_NAME"),
			Port:     v.GetString("DB_PORT"),

			SSLMode:  v.GetString("SSL_MODE"),
			TimeZone: v.GetString("TIME_ZONE"),

			MaxIdleConns:    v.GetInt("MAX_IDLE_CONNS"),
			MaxOpenConns:    v.GetInt("MAX_OPEN_CONNS"),
			ConnMaxLifetime: v.GetDuration("CONN_MAX_LIFETIME"),
		},
		JwtConfig: JwtConfig{
			AccessTTL:  v.GetDuration("JWT_ACCESS_TTL"),
			RefreshTTL: v.GetDuration("JWT_REFRESH_TTL"),
		},
		RedisConfig: RedisConfig{
			Addr:     v.GetString("REDIS_ADDR"),
			Password: v.GetString("REDIS_PASSWORD"),
			DB:       v.GetInt("REDIS_DB"),
		},
		MinioConfig: MinioConfig{
			Endpoint:        v.GetString("MINIO_ENDPOINT"),
			AccessKeyID:     v.GetString("MINIO_ACCESS_KEY_ID"),
			SecretAccessKey: v.GetString("MINIO_SECRET_ACCESS_KEY"),
			UseSSL:          v.GetBool("MINIO_USE_SSL"),
			Bucket:          v.GetString("MINIO_BUCKET"),
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

	v.SetDefault("REDIS_ADDR", "127.0.0.1:6379")
	v.SetDefault("REDIS_PASSWORD", "")
	v.SetDefault("REDIS_DB", 0)

	v.SetDefault("MINIO_ENDPOINT", "127.0.0.1:9000")
	v.SetDefault("MINIO_USE_SSL", false)
	v.SetDefault("MINIO_BUCKET", "usercontent")
}

func validateConfig(cfg *Config) error {
	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return fmt.Errorf("missing required attributes: %w", err)
	}
	return nil
}
