package main

import (
	"blog-api/config"
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
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

func main() {
	if err := config.Init(); err != nil {
		log.Fatalf("init config error: %v", err)
	}

	cfg := config.GetConfig()

	if err := database.InitDB(database.BuildDSN(cfg.DbHost, cfg.DbUser, cfg.DbPassword, cfg.DbName, cfg.DbPort), &database.PoolConfig{
		MaxIdleConns:    cfg.MaxIdleConns,
		MaxOpenConns:    cfg.MaxOpenConns,
		ConnMaxLifetime: cfg.ConnMaxLifetime,
	}); err != nil {
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

	jwtManager := jwtmanager.NewJWTManager(cfg.SecretKey, cfg.JwtAccessTTL, cfg.JwtRefreshTTL)
	userService := users.NewUserService(jwtManager)
	middlewareManager := middleware.NewManager(jwtManager, userService)

	apiGroup := app.Group("/api")
	{
		authGroup := apiGroup.Group("/auth")
		authService := auth.NewAuthService(jwtManager)
		authHandler := auth.NewAuthHandler(authService)
		routes.RegisterAuthRoutes(authGroup, authHandler)
	}

	{
		usersGroup := apiGroup.Group("/users")
		userHandler := users.NewUserHandler()
		routes.RegisterUserRoutes(usersGroup, userHandler, middlewareManager)
	}

	{
		postsGroup := apiGroup.Group("/posts")
		postService := posts.NewPostService()
		postHandler := posts.NewPostHandler(postService)
		routes.RegisterPostRoutes(postsGroup, postHandler, middlewareManager)
	}

	log.Fatal(app.Listen(":3000"))
}
