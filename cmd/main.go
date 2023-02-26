package main

import (
	"github.com/devararishivian/antrekuy/internal/config"
	"github.com/devararishivian/antrekuy/internal/infrastructure"
	"github.com/devararishivian/antrekuy/internal/interface/api"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"
)

func main() {
	errConfig := config.LoadConfig("./internal/config/config.json")
	if errConfig != nil {
		panic(errConfig)
	}

	db, errDB := infrastructure.NewDatabase()
	if errDB != nil {
		panic(errDB)
	}

	app := fiber.New()
	app.Use(logger.New())

	api.RegisterRoutes(app.Group("api"), db)

	log.Fatal(app.Listen(config.Server.Address))
}
