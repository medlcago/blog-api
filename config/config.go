package config

import (
	"fmt"
	"sync"

	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	SecretKey string `validate:"required"`

	DbHost     string `validate:"required"`
	DbUser     string `validate:"required"`
	DbPassword string `validate:"required"`
	DbName     string `validate:"required"`
	DbPort     string `validate:"required"`

	MaxIdleConns    int           `validate:"required"`
	MaxOpenConns    int           `validate:"required"`
	ConnMaxLifetime time.Duration `validate:"required"`

	JwtAccessTTL  time.Duration `validate:"required"`
	JwtRefreshTTL time.Duration `validate:"required"`
}

var (
	configInstance *Config
	configOnce     sync.Once
	configErr      error
)

func Init() error {
	configOnce.Do(func() {
		configInstance, configErr = loadConfig()
	})
	return configErr
}

func GetConfig() *Config {
	if configInstance == nil {
		panic("config not initialized. Call Init() first")
	}
	return configInstance
}

func loadConfig() (*Config, error) {
	v := viper.New()
	v.SetConfigFile(".env")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	v.AutomaticEnv()

	setDefaults(v)

	config := &Config{
		SecretKey:       v.GetString("SECRET_KEY"),
		DbHost:          v.GetString("DB_HOST"),
		DbUser:          v.GetString("DB_USER"),
		DbPassword:      v.GetString("DB_PASSWORD"),
		DbName:          v.GetString("DB_NAME"),
		DbPort:          v.GetString("DB_PORT"),
		MaxIdleConns:    v.GetInt("MAX_IDLE_CONNS"),
		MaxOpenConns:    v.GetInt("MAX_OPEN_CONNS"),
		ConnMaxLifetime: v.GetDuration("CONN_MAX_LIFETIME"),
		JwtAccessTTL:    v.GetDuration("JWT_ACCESS_TTL"),
		JwtRefreshTTL:   v.GetDuration("JWT_REFRESH_TTL"),
	}

	if err := validateConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

func setDefaults(v *viper.Viper) {
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
