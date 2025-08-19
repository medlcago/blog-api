package main

import (
	"blog-api/internal/auth"
	"blog-api/internal/database"
	"blog-api/internal/middleware"
	"blog-api/internal/posts"
	"blog-api/internal/routes"
	"blog-api/internal/users"
	appvalidator "blog-api/internal/validator"
	"blog-api/pkg/errors"
	"blog-api/pkg/jwtmanager"
	"blog-api/pkg/struct_validator"
	"fmt"
	"log"
	"os"
	"time"

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

	validator, err := appvalidator.New()
	if err != nil {
		log.Fatalf("validator error: %v", err)
	}

	app := fiber.New(fiber.Config{
		StructValidator: struct_validator.New(validator),
		ErrorHandler:    errors.ErrorHandler,
	})

	app.Use(logger.New())
	app.Use(recover.New())

	jwtManager := jwtmanager.NewJWTManager(os.Getenv("SECRET_KEY"), 30*time.Minute, 24*time.Hour)
	userService := users.NewUserService(jwtManager)
	middlewareManager := middleware.NewManager(jwtManager, userService)

	{
		authGroup := app.Group("/auth")
		authService := auth.NewAuthService(jwtManager)
		authHandler := auth.NewAuthHandler(authService)
		routes.RegisterAuthRoutes(authGroup, authHandler)
	}

	{
		usersGroup := app.Group("/users")
		userHandler := users.NewUserHandler()
		routes.RegisterUserRoutes(usersGroup, userHandler, middlewareManager)
	}

	{
		postsGroup := app.Group("/posts")
		postService := posts.NewPostService()
		postHandler := posts.NewPostHandler(postService)
		routes.RegisterPostRoutes(postsGroup, postHandler, middlewareManager)
	}

	log.Fatal(app.Listen(":3000"))
}
