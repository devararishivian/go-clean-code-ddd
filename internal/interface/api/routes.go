package api

import (
	"github.com/devararishivian/antrekuy/internal/application/usecase"
	"github.com/devararishivian/antrekuy/internal/infrastructure"
	"github.com/devararishivian/antrekuy/internal/infrastructure/memory"
	"github.com/devararishivian/antrekuy/internal/infrastructure/persistence"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(router fiber.Router, db *infrastructure.Database, redisClient *infrastructure.Redis) {
	router.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	v1Route := router.Group("v1")

	registerUserRoutesV1(v1Route, db)
	registerAuthRoutesV1(v1Route, db, redisClient)
}

func registerUserRoutesV1(router fiber.Router, db *infrastructure.Database) {
	repository := persistence.NewUserRepository(db)
	useCase := usecase.NewUserUseCase(repository)
	handler := NewUserHandler(useCase)

	route := router.Group("user")
	route.Post("/", handler.Store)
}

func registerAuthRoutesV1(router fiber.Router, db *infrastructure.Database, redisClient *infrastructure.Redis) {
	userRepository := persistence.NewUserRepository(db)
	userUseCase := usecase.NewUserUseCase(userRepository)

	cacheRepository := memory.NewCacheRepository(redisClient)
	useCase := usecase.NewAuthUseCase(userUseCase, cacheRepository)
	handler := NewAuthHandler(useCase)

	route := router.Group("auth")
	route.Post("/", handler.Authenticate)
	route.Post("/refresh", handler.RefreshToken)
}
