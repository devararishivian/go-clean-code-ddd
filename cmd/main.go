package main

import (
	"github.com/devararishivian/antrekuy/internal/application/usecase"
	"github.com/devararishivian/antrekuy/internal/infrastructure"
	"github.com/devararishivian/antrekuy/internal/infrastructure/persistence"
	"github.com/devararishivian/antrekuy/internal/interface/api"
	"github.com/gofiber/fiber/v2"
	"log"
)

func main() {
	// Create a new database connection
	db, err := infrastructure.NewDatabase()
	if err != nil {
		panic(err)
	}

	userRepository := persistence.NewUserRepository(db)
	userUseCase := usecase.NewUserUseCase(userRepository)

	app := fiber.New()
	apiPrefixV1 := app.Group("api/v1")

	api.RegisterUserRoutes(apiPrefixV1, userUseCase)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	log.Fatal(app.Listen(":3000"))
}
