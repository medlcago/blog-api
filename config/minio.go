package config

import "github.com/spf13/viper"

type MinioConfig struct {
	Endpoint        string `validate:"required"`
	AccessKeyID     string `validate:"required"`
	SecretAccessKey string `validate:"required"`
	UseSSL          bool
	Bucket          string `validate:"required"`
}

func loadMinioConfig(v *viper.Viper) MinioConfig {
	return MinioConfig{
		Endpoint:        v.GetString("MINIO_ENDPOINT"),
		AccessKeyID:     v.GetString("MINIO_ACCESS_KEY_ID"),
		SecretAccessKey: v.GetString("MINIO_SECRET_ACCESS_KEY"),
		UseSSL:          v.GetBool("MINIO_USE_SSL"),
		Bucket:          v.GetString("MINIO_BUCKET"),
	}
}
