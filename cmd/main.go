package main

import (
	"blog-api/config"
	"blog-api/internal/database"
	"blog-api/internal/server"
	"blog-api/internal/storage"
	appvalidator "blog-api/internal/validator"
	"log"
	"os"
	"time"
)

func init() {
	time.Local = time.UTC
}

func main() {
	appLogger := log.New(os.Stderr, "APP: ", log.Ldate|log.Ltime|log.Lshortfile)

	appLogger.Println("Starting application initialization...")

	cfg := config.MustGet()

	redisClient, err := storage.NewRedisClient(cfg.RedisConfig)
	if err != nil {
		appLogger.Fatalf("failed to init redis client: %v", err)
	}

	minioClient, err := storage.NewMinioClient(cfg.MinioConfig)
	if err != nil {
		appLogger.Fatalf("failed to create minio client: %v", err)
	}

	db, err := database.New(cfg.DatabaseConfig)
	if err != nil {
		appLogger.Fatalf("failed to init database: %v", err)
	}

	if err = db.RunMigrations(); err != nil {
		appLogger.Fatalf("failed to run database migrations: %v", err)
	}
	appLogger.Println("✅ Database migrations completed")

	validate, err := appvalidator.New()
	if err != nil {
		appLogger.Fatalf("failed to init validator: %v", err)
	}

	serverDeps := &server.Dependencies{
		Cfg:         cfg,
		DB:          db,
		RedisClient: redisClient,
		MinioClient: minioClient,
		Validate:    validate,
		AppLogger:   appLogger,
	}
	s, err := server.NewServer(serverDeps)
	if err != nil {
		appLogger.Fatalf("failed to init server: %v", err)
	}

	if err := s.Run(); err != nil {
		appLogger.Fatalf("server stopped with error: %v", err)
	}
}
