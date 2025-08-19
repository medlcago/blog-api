package middleware

import (
	"blog-api/internal/users"
	"blog-api/pkg/jwtmanager"
)

type Manager struct {
	jwtManager  *jwtmanager.JWTManager
	userService users.IUserService
}

func NewManager(jwtManager *jwtmanager.JWTManager, userService users.IUserService) *Manager {
	return &Manager{
		jwtManager:  jwtManager,
		userService: userService,
	}
}
