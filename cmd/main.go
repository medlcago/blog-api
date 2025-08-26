package main

import (
	"blog-api/config"
	"blog-api/internal/database"
	"blog-api/internal/server"
	appvalidator "blog-api/internal/validator"
	redisStore "blog-api/pkg/storage/redis"
	"context"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

func init() {
	time.Local = time.UTC
}

func main() {
	appLogger := log.New(os.Stderr, "APP: ", log.Ldate|log.Ltime|log.Lshortfile)

	appLogger.Println("Starting application initialization...")

	cfg := config.MustGet()

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisConfig.Addr,
		Password: cfg.RedisConfig.Password,
		DB:       cfg.RedisConfig.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		appLogger.Fatalf("failed to connect to redis: %v", err)
	}

	store, err := redisStore.New(rdb, nil)
	if err != nil {
		appLogger.Fatalf("failed to init store: %v", err)
	}

	db, err := database.New(cfg.DatabaseConfig)
	if err != nil {
		appLogger.Fatalf("failed to init database: %v", err)
	}

	if err = db.RunMigrations(); err != nil {
		appLogger.Fatalf("failed to run storage migrations: %v", err)
	}
	appLogger.Println("✅ Database migrations completed")

	validate, err := appvalidator.New()
	if err != nil {
		appLogger.Fatalf("failed to init validator: %v", err)
	}

	serverDeps := &server.Dependencies{
		Cfg:       cfg,
		DB:        db,
		Store:     store,
		Validate:  validate,
		AppLogger: appLogger,
	}
	s, err := server.NewServer(serverDeps)
	if err != nil {
		appLogger.Fatalf("failed to init server: %v", err)
	}

	if err := s.Run(); err != nil {
		appLogger.Fatalf("server stopped with error: %v", err)
	}
}
