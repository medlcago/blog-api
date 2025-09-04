package middleware

import (
	"blog-api/internal/tokenmanager"
	"blog-api/internal/users"
	"log/slog"
)

type Manager struct {
	log         *slog.Logger
	jwtService  tokenmanager.TokenManager
	userService users.IUserService
}

func NewManager(log *slog.Logger, jwtService tokenmanager.TokenManager, userService users.IUserService) *Manager {
	return &Manager{
		log:         log,
		jwtService:  jwtService,
		userService: userService,
	}
}
