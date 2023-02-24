package api

import (
	"github.com/devararishivian/antrekuy/internal/domain/service"
	"github.com/gofiber/fiber/v2"
)

// RegisterUserRoutes registers user-related routes
func RegisterUserRoutes(router fiber.Router, userUseCase service.UserService) {
	handler := NewUserHandler(userUseCase)

	router.Get("/user", handler.Store)
}
