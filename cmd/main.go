package main

import (
	"blog-api/config"
	"blog-api/internal/database"
	"blog-api/internal/server"
	appvalidator "blog-api/internal/validator"
	"log"
	"os"
)

func main() {
	logger := log.New(os.Stderr, "APP: ", log.Ldate|log.Ltime|log.Lshortfile)

	logger.Println("Starting application initialization...")

	cfg := config.MustGet()
	db, err := database.New(cfg.DatabaseConfig)
	if err != nil {
		logger.Fatalf("failed to init database: %v", err)
	}

	if err = db.RunMigrations(); err != nil {
		logger.Fatalf("failed to run database migrations: %v", err)
	}
	logger.Println("✅ Database migrations completed")

	validate, err := appvalidator.New()
	if err != nil {
		logger.Fatalf("failed to init validator: %v", err)
	}

	s, err := server.NewServer(cfg, db, validate, logger)
	if err != nil {
		logger.Fatalf("failed to init server: %v", err)
	}

	if err := s.Run(); err != nil {
		logger.Fatalf("server stopped with error: %v", err)
	}
}
