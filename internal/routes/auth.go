package routes

import (
	"blog-api/internal/auth"
	"blog-api/internal/middleware"

	"github.com/gofiber/fiber/v3"
)

func RegisterAuthRoutes(r fiber.Router, h auth.IAuthHandler, mw *middleware.Manager) {
	r.Post("/register", h.Register)
	r.Post("/login", h.Login)
	r.Post("/login/2fa", h.Login2FA)
	r.Post("/refresh-token", h.RefreshToken)
	r.Post("/change-password", mw.AuthMiddleware(), h.ChangePassword)

	r.Post("/enable-2fa", mw.AuthMiddleware(), h.Enable2FA)
	r.Post("/verify-2fa", mw.AuthMiddleware(), h.Verify2FA)
	r.Post("/disable-2fa", mw.AuthMiddleware(), h.Disable2FA)
}
