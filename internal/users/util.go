package users

import "github.com/gofiber/fiber/v3"

func GetUser(ctx fiber.Ctx) *UserResponse {
	user := fiber.Locals[*UserResponse](ctx, "user")
	return user
}

func MustGetUser(ctx fiber.Ctx) *UserResponse {
	user := GetUser(ctx)
	if user == nil {
		panic("user not found in context (did you forget AuthMiddleware?)")
	}
	return user
}
