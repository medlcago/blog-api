package config

import "github.com/spf13/viper"

type RedisConfig struct {
	Addr     string `validate:"required"`
	Password string
	DB       int
}

func loadRedisConfig(v *viper.Viper) RedisConfig {
	return RedisConfig{
		Addr:     v.GetString("REDIS_ADDR"),
		Password: v.GetString("REDIS_PASSWORD"),
		DB:       v.GetInt("REDIS_DB"),
	}
}
