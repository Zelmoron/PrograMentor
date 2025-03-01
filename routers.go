package main

import (
	"github.com/gofiber/fiber/v2"
	"main/handlers"
	"os"
)

func initRoutes(app *fiber.App, in *handlers.InHandlers, out *handlers.OutHandlers) {
	app.Get("/", func(c *fiber.Ctx) error { return c.JSON(os.Getenv("APP_VERSION")) })

}
