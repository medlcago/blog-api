package config

type RedisConfig struct {
	Addr     string `validate:"required"`
	Password string
	DB       int
}
