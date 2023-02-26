package main

import (
	"github.com/devararishivian/antrekuy/internal/infrastructure"
	"github.com/devararishivian/antrekuy/internal/interface/api"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"
)

func main() {
	db, err := infrastructure.NewDatabase()
	if err != nil {
		panic(err)
	}

	app := fiber.New()
	app.Use(logger.New())

	api.RegisterRoutes(app.Group("api"), db)

	log.Fatal(app.Listen(":3000"))
}
