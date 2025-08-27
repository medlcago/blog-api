package config

type MinioConfig struct {
	Endpoint        string `validate:"required"`
	AccessKeyID     string `validate:"required"`
	SecretAccessKey string `validate:"required"`
	UseSSL          bool
	Bucket          string `validate:"required"`
}
