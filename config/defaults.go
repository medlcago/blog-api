package config

import (
	"time"

	"github.com/spf13/viper"
)

func setDefaults(v *viper.Viper) {
	v.SetDefault("APP_ENV", "dev")

	v.SetDefault("SERVER_HOST", "0.0.0.0")
	v.SetDefault("SERVER_PORT", "3000")
	v.SetDefault("SERVER_BODY_LIMIT", 10<<20)
	v.SetDefault("SERVER_READ_TIMEOUT", 5*time.Second)
	v.SetDefault("SERVER_WRITE_TIMEOUT", 5*time.Second)
	v.SetDefault("SERVER_SHUTDOWN_TIMEOUT", 30*time.Second)

	v.SetDefault("DB_SSL_MODE", "disable")
	v.SetDefault("DB_TIME_ZONE", "UTC")

	v.SetDefault("DB_MAX_IDLE_CONNS", 10)
	v.SetDefault("DB_MAX_OPEN_CONNS", 100)
	v.SetDefault("DB_CONN_MAX_LIFETIME", time.Hour)

	v.SetDefault("JWT_ACCESS_TTL", 30*time.Minute)
	v.SetDefault("JWT_REFRESH_TTL", 24*30*time.Hour)

	v.SetDefault("REDIS_ADDR", "127.0.0.1:6379")
	v.SetDefault("REDIS_PASSWORD", "")
	v.SetDefault("REDIS_DB", 0)

	v.SetDefault("MINIO_ENDPOINT", "127.0.0.1:9000")
	v.SetDefault("MINIO_USE_SSL", false)
	v.SetDefault("MINIO_BUCKET", "usercontent")
}
