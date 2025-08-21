package config

type ServerConfig struct {
	ServerHost string `validate:"required"`
	ServerPort string `validate:"required"`
}
