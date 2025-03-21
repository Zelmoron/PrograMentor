package main

import (
	"os"

	"github.com/gofiber/fiber/v2"

	"main/handlers"
)

func initRoutes(app *fiber.App, in *handlers.InHandlers, out *handlers.OutHandlers) {
	app.Get("/", func(c *fiber.Ctx) error { return c.JSON(os.Getenv("APP_VERSION")) })

	// Без аутентификации
	api := app.Group("")
	api.Post("/auth", in.Login, out.LoginOut)
	api.Post("/refresh-token", out.RefreshToken)

	protected := api.Group("/protected", handlers.JWT())
	protected.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello, authenticated user",
		})
	})
}
