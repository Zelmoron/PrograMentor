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
	api.Post("/auth", in.Login)
	api.Post("/refresh-token", out.RefreshToken)

	protected := api.Group("/protected")

	// Обработчик для защищенного маршрута
	protected.Get("/", handlers.JWT(out.GetUsers()), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello, authenticated user",
		})
	})
}
