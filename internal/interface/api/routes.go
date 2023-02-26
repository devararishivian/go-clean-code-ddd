package api

import (
	"github.com/devararishivian/antrekuy/internal/application/usecase"
	"github.com/devararishivian/antrekuy/internal/infrastructure"
	"github.com/devararishivian/antrekuy/internal/infrastructure/persistence"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(router fiber.Router, db *infrastructure.Database) {
	router.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	v1Route := router.Group("v1")

	registerUserRoutesV1(v1Route, db)
}

func registerUserRoutesV1(router fiber.Router, db *infrastructure.Database) {
	repository := persistence.NewUserRepository(db)
	useCase := usecase.NewUserUseCase(repository)
	handler := NewUserHandler(useCase)

	userRoute := router.Group("user")
	userRoute.Get("/", handler.Store)
}
