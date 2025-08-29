package main

import (
	"blog-api/config"
	"blog-api/internal/database"
	"blog-api/internal/logger"
	"blog-api/internal/server"
	"blog-api/internal/storage"
	appvalidator "blog-api/internal/validator"
	"log/slog"
	"os"
	"time"
)

func init() {
	time.Local = time.UTC
}

func main() {
	cfg := config.MustGet()

	log := logger.New(logger.Env(cfg.Env))

	log.Info("Starting application initialization...")

	redisClient, err := storage.NewRedisClient(cfg.RedisConfig)
	if err != nil {
		log.Error("failed to init redis client", slog.Any("error", err))
		os.Exit(1)
	}

	minioClient, err := storage.NewMinioClient(cfg.MinioConfig)
	if err != nil {
		log.Error("failed to create minio client", slog.Any("error", err))
		os.Exit(1)
	}

	db, err := database.New(cfg.DatabaseConfig)
	if err != nil {
		log.Error("failed to init database", slog.Any("error", err))
		os.Exit(1)
	}

	if err = db.RunMigrations(); err != nil {
		log.Error("failed to run database migrations", slog.Any("error", err))
		os.Exit(1)
	}
	log.Info("✅ Database migrations completed")

	validator, err := appvalidator.New()
	if err != nil {
		log.Error("failed to init validator", slog.Any("error", err))
		os.Exit(1)
	}

	serverDeps := &server.Dependencies{
		Cfg:         cfg,
		DB:          db,
		RedisClient: redisClient,
		MinioClient: minioClient,
		Validator:   validator,
		Logger:      log,
	}
	s, err := server.NewServer(serverDeps)
	if err != nil {
		log.Error("failed to init server", slog.Any("error", err))
		os.Exit(1)
	}

	if err := s.Run(); err != nil {
		log.Error("server stopped with error", slog.Any("error", err))
		os.Exit(1)
	}
}
