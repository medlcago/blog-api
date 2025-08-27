package server

import (
	"blog-api/config"
	"blog-api/internal/auth"
	"blog-api/internal/database"
	"blog-api/internal/jwtmanager"
	"blog-api/internal/middleware"
	"blog-api/internal/photos"
	"blog-api/internal/posts"
	"blog-api/internal/routes"
	internalStorage "blog-api/internal/storage"
	"blog-api/internal/users"
	"blog-api/pkg/errors"
	"blog-api/pkg/storage"
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

type Dependencies struct {
	Cfg         *config.Config
	DB          *database.DB
	Store       storage.Storage
	MinioClient *internalStorage.MinioClient
	Validate    *validator.Validate
	AppLogger   *log.Logger
}

type Server struct {
	app *fiber.App
	*Dependencies
}

func NewServer(deps *Dependencies) (*Server, error) {
	if deps == nil {
		panic("NewServer: deps cannot be nil")
	}

	app := fiber.New(fiber.Config{
		StructValidator: struct_validator.New(deps.Validate),
		ErrorHandler:    errors.ErrorHandler,
	})

	app.Use(logger.New())
	app.Use(recover.New())

	jwtManager := jwtmanager.NewJWTManager(deps.Cfg.SecretKey, deps.Cfg.JwtConfig)

	photoProcessor := photos.NewProcessor(5*1024*1024, []string{"jpeg", "png"})

	// Services
	userService := users.NewUserService(jwtManager, deps.DB)
	authService := auth.NewAuthService(jwtManager, deps.DB, deps.Store.WithNamespace("authService"), deps.AppLogger)
	postService := posts.NewPostService(deps.DB)
	photoService := photos.NewPhotoService(deps.DB, deps.MinioClient, photoProcessor)

	middlewareManager := middleware.NewManager(jwtManager, userService)

	// Handlers
	authHandler := auth.NewAuthHandler(authService)
	userHandler := users.NewUserHandler(photoService)
	postHandler := posts.NewPostHandler(postService)

	// Groups
	apiGroup := app.Group("/api")
	authGroup := apiGroup.Group("/auth")
	usersGroup := apiGroup.Group("/users")
	postsGroup := apiGroup.Group("/posts")

	// Routes
	routes.RegisterAuthRoutes(authGroup, authHandler)
	routes.RegisterUserRoutes(usersGroup, userHandler, middlewareManager)
	routes.RegisterPostRoutes(postsGroup, postHandler, middlewareManager)

	return &Server{
		app:          app,
		Dependencies: deps,
	}, nil
}

func (s *Server) Run() error {
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	serverErr := make(chan error, 1)
	go func() {
		addr := fmt.Sprintf("%s:%s", s.Cfg.ServerConfig.Host, s.Cfg.ServerConfig.Port)
		s.AppLogger.Printf("Server starting on %s", addr)
		if err := s.app.Listen(addr); err != nil {
			serverErr <- err
		}
	}()

	select {
	case err1 := <-serverErr:
		s.AppLogger.Printf("Server error: %v", err1)
		err2 := s.DB.Close()
		return goerrors.Join(err1, err2)

	case sig := <-shutdownChan:
		s.AppLogger.Printf("Received signal: %v. Shutting down gracefully...", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := s.app.ShutdownWithContext(ctx); err != nil {
			if goerrors.Is(err, context.DeadlineExceeded) {
				s.AppLogger.Printf("Shutdown timed out after 30 seconds, forcing exit")
			} else {
				s.AppLogger.Printf("Error during shutdown: %v", err)
			}
		} else {
			s.AppLogger.Println("Server stopped gracefully")
		}

		if err := s.DB.Close(); err != nil {
			s.AppLogger.Printf("Database close error: %v", err)
		}

		s.AppLogger.Println("Shutdown completed")
		return nil
	}
}
