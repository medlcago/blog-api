package config

type ServerConfig struct {
	Host string `validate:"required"`
	Port string `validate:"required"`
}
