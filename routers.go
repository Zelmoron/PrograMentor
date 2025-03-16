package main

import (
	"main/handlers"
	"os"

	"github.com/gofiber/fiber/v2"
)

func initRoutes(app *fiber.App, in *handlers.InHandlers, out *handlers.OutHandlers) {
	app.Get("/", func(c *fiber.Ctx) error { return c.JSON(os.Getenv("APP_VERSION")) })

	//Without autification
	api := app.Group("")

	api.Post("/auth", in.Login)

}
