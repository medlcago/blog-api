package middleware

import (
	"blog-api/internal/jwtmanager"
	"blog-api/internal/users"
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
