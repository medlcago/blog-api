package main

import (
	"blog-api/internal/database"
	"blog-api/internal/posts"
	"blog-api/pkg/errors"
	structValidator "blog-api/pkg/validator"
	"fmt"
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	if err := database.InitDB(dsn); err != nil {
		log.Fatalf("database init error: %v", err)
	}

	if err := database.RunMigrations(); err != nil {
		log.Fatalf("database migration error: %v", err)
	}

	app := fiber.New(fiber.Config{
		StructValidator: structValidator.New(validator.New()),
		ErrorHandler:    errors.ErrorHandler,
	})

	app.Use(logger.New())
	app.Use(recover.New())

	{
		postsGroup := app.Group("/posts")
		postService := posts.NewPostService()
		postHandler := posts.NewPostHandler(postService)
		posts.RegisterRoutes(postsGroup, postHandler)
	}

	log.Fatal(app.Listen(":3000"))
}
