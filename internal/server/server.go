package server

import (
	"blog-api/config"
	"blog-api/internal/auth"
	"blog-api/internal/database"
	"blog-api/internal/middleware"
	"blog-api/internal/posts"
	"blog-api/internal/routes"
	"blog-api/internal/users"
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

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

type Server struct {
	app       *fiber.App
	cfg       *config.Config
	db        *database.DB
	validate  *validator.Validate
	appLogger *log.Logger
}

func NewServer(cfg *config.Config, db *database.DB, validate *validator.Validate, appLogger *log.Logger) (*Server, error) {
	app := fiber.New(fiber.Config{
		StructValidator: struct_validator.New(validate),
		ErrorHandler:    errors.ErrorHandler,
	})

	app.Use(logger.New())
	app.Use(recover.New())

	jwtManager := jwtmanager.NewJWTManager(cfg.SecretKey, cfg.JwtConfig)
	userService := users.NewUserService(jwtManager, db)
	middlewareManager := middleware.NewManager(jwtManager, userService)

	apiGroup := app.Group("/api")

	// Auth
	authGroup := apiGroup.Group("/auth")
	authService := auth.NewAuthService(jwtManager, db)
	authHandler := auth.NewAuthHandler(authService)
	routes.RegisterAuthRoutes(authGroup, authHandler)

	// Users
	usersGroup := apiGroup.Group("/users")
	userHandler := users.NewUserHandler()
	routes.RegisterUserRoutes(usersGroup, userHandler, middlewareManager)

	// Posts
	postsGroup := apiGroup.Group("/posts")
	postService := posts.NewPostService(db)
	postHandler := posts.NewPostHandler(postService)
	routes.RegisterPostRoutes(postsGroup, postHandler, middlewareManager)

	return &Server{
		app:       app,
		cfg:       cfg,
		db:        db,
		validate:  validate,
		appLogger: appLogger,
	}, nil
}

func (s *Server) Run() error {
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	serverErr := make(chan error, 1)
	go func() {
		addr := fmt.Sprintf("%s:%s", s.cfg.ServerConfig.ServerHost, s.cfg.ServerConfig.ServerPort)
		s.appLogger.Printf("Server starting on %s", addr)
		if err := s.app.Listen(addr); err != nil {
			serverErr <- err
		}
	}()

	select {
	case err1 := <-serverErr:
		s.appLogger.Printf("Server error: %v", err1)
		err2 := s.db.Close()
		return goerrors.Join(err1, err2)

	case sig := <-shutdownChan:
		s.appLogger.Printf("Received signal: %v. Shutting down gracefully...", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := s.app.ShutdownWithContext(ctx); err != nil {
			if goerrors.Is(err, context.DeadlineExceeded) {
				s.appLogger.Printf("Shutdown timed out after 30 seconds, forcing exit")
			} else {
				s.appLogger.Printf("Error during shutdown: %v", err)
			}
		} else {
			s.appLogger.Println("Server stopped gracefully")
		}

		if err := s.db.Close(); err != nil {
			s.appLogger.Printf("Database close error: %v", err)
		}

		s.appLogger.Println("Shutdown completed")
		return nil
	}
}
