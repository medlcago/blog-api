package server

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
	"context"
	goerrors "errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

type Server struct {
	app *fiber.App
	cfg *config.Config
}

func NewServer(cfg *config.Config) (*Server, error) {
	if err := database.InitDB(
		database.BuildDSN(cfg.DbHost, cfg.DbUser, cfg.DbPassword, cfg.DbName, cfg.DbPort),
		&database.PoolConfig{
			MaxIdleConns:    cfg.MaxIdleConns,
			MaxOpenConns:    cfg.MaxOpenConns,
			ConnMaxLifetime: cfg.ConnMaxLifetime,
		},
	); err != nil {
		return nil, err
	}
	if err := database.RunMigrations(); err != nil {
		return nil, err
	}

	validator, err := appvalidator.New()
	if err != nil {
		return nil, err
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

	// Auth
	authGroup := apiGroup.Group("/auth")
	authService := auth.NewAuthService(jwtManager)
	authHandler := auth.NewAuthHandler(authService)
	routes.RegisterAuthRoutes(authGroup, authHandler)

	// Users
	usersGroup := apiGroup.Group("/users")
	userHandler := users.NewUserHandler()
	routes.RegisterUserRoutes(usersGroup, userHandler, middlewareManager)

	// Posts
	postsGroup := apiGroup.Group("/posts")
	postService := posts.NewPostService()
	postHandler := posts.NewPostHandler(postService)
	routes.RegisterPostRoutes(postsGroup, postHandler, middlewareManager)

	return &Server{
		app: app,
		cfg: cfg,
	}, nil
}

func (s *Server) Run() error {
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	serverErr := make(chan error, 1)
	go func() {
		addr := fmt.Sprintf("%s:%s", s.cfg.ServerHost, s.cfg.ServerPort)
		log.Printf("Server starting on %s", addr)
		if err := s.app.Listen(addr); err != nil {
			serverErr <- err
		}
	}()

	select {
	case err1 := <-serverErr:
		log.Printf("Server error: %v", err1)
		err2 := database.Close()
		return goerrors.Join(err1, err2)

	case sig := <-shutdownChan:
		log.Printf("Received signal: %v. Shutting down gracefully...", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := s.app.ShutdownWithContext(ctx); err != nil {
			if goerrors.Is(err, context.DeadlineExceeded) {
				log.Printf("Shutdown timed out after 30 seconds, forcing exit")
			} else {
				log.Printf("Error during shutdown: %v", err)
			}
		} else {
			log.Println("Server stopped gracefully")
		}

		if err := database.Close(); err != nil {
			log.Printf("Database close error: %v", err)
		}

		log.Println("Shutdown completed")
		return nil
	}
}
