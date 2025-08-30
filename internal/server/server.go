package server

import (
	"blog-api/config"
	"blog-api/internal/auth"
	"blog-api/internal/database"
	"blog-api/internal/errors"
	"blog-api/internal/middleware"
	"blog-api/internal/photos"
	"blog-api/internal/posts"
	"blog-api/internal/routes"
	"blog-api/internal/storage"
	"blog-api/internal/tokenmanager"
	"blog-api/internal/users"
	"blog-api/pkg/struct_validator"
	"context"
	goerrors "errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
)

type Dependencies struct {
	Cfg         *config.Config
	DB          *database.DB
	RedisClient *storage.RedisClient
	MinioClient *storage.MinioClient
	Validator   *validator.Validate
	Logger      *slog.Logger
}

type Server struct {
	app *fiber.App
	*Dependencies
}

func NewServer(deps *Dependencies) (*Server, error) {
	if deps == nil {
		panic("NewServer: deps cannot be nil")
	}

	// Services
	jwtService := tokenmanager.NewJWTManager(deps.Cfg.SecretKey, deps.Cfg.JwtConfig)
	userService := users.NewUserService(deps.DB, deps.Logger)
	authService := auth.NewAuthService(jwtService, deps.DB, deps.RedisClient, deps.Logger)
	postService := posts.NewPostService(deps.DB, deps.Logger)
	photoService := photos.NewPhotoService(deps.DB, deps.MinioClient, deps.Logger)

	mw := middleware.NewManager(deps.Logger, jwtService, userService)

	// Handlers
	authHandler := auth.NewAuthHandler(authService)
	userHandler := users.NewUserHandler(photoService)
	postHandler := posts.NewPostHandler(postService)

	// App
	app := fiber.New(fiber.Config{
		StructValidator: struct_validator.New(deps.Validator),
		ErrorHandler:    errors.NewErrorHandler(deps.Logger),
		BodyLimit:       deps.Cfg.ServerConfig.BodyLimit,
		ReadTimeout:     deps.Cfg.ServerConfig.ReadTimeout,
		WriteTimeout:    deps.Cfg.ServerConfig.WriteTimeout,
	})

	app.Use(requestid.New())
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))
	app.Use(mw.LoggerMiddleware())

	// Groups
	apiGroup := app.Group("/api")
	authGroup := apiGroup.Group("/auth")
	usersGroup := apiGroup.Group("/users")
	postsGroup := apiGroup.Group("/posts")

	// Routes
	routes.RegisterAuthRoutes(authGroup, authHandler)
	routes.RegisterUserRoutes(usersGroup, userHandler, mw)
	routes.RegisterPostRoutes(postsGroup, postHandler, mw)

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
		s.Logger.Info(
			"starting server",
			slog.String("addr", addr),
			slog.String("env", s.Cfg.Env),
		)
		if err := s.app.Listen(addr); err != nil {
			serverErr <- err
		}
	}()

	select {
	case err1 := <-serverErr:
		s.Logger.Error("Server error", slog.Any("error", err1))
		err2 := s.DB.Close()
		return goerrors.Join(err1, err2)

	case sig := <-shutdownChan:
		s.Logger.Info("Received signal. Shutting down gracefully...", slog.String("signal", sig.String()))

		ctx, cancel := context.WithTimeout(context.Background(), s.Cfg.ServerConfig.ShutdownTimeout)
		defer cancel()

		if err := s.app.ShutdownWithContext(ctx); err != nil {
			if goerrors.Is(err, context.DeadlineExceeded) {
				s.Logger.Info(
					"Shutdown timed out, forcing exit",
					slog.Duration("timeout", s.Cfg.ServerConfig.ShutdownTimeout),
				)
			} else {
				s.Logger.Info("Error during shutdown", slog.Any("error", err))
			}
		} else {
			s.Logger.Info("Server stopped gracefully")
		}

		if err := s.DB.Close(); err != nil {
			s.Logger.Error("Database close error", slog.Any("error", err))
		}

		if err := s.RedisClient.Client.Close(); err != nil {
			s.Logger.Error("Redis close error", slog.Any("error", err))
		}

		s.Logger.Info("Shutdown completed")
		return nil
	}
}
