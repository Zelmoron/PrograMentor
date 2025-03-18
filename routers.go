package main

import (
	"os"

	"github.com/gofiber/fiber/v2"

	"main/handlers"
)

func initRoutes(app *fiber.App, in *handlers.InHandlers, out *handlers.OutHandlers) {
	app.Get("/", func(c *fiber.Ctx) error { return c.JSON(os.Getenv("APP_VERSION")) })

	//Without autification
	api := app.Group("")
	api.Post("/auth", in.Login)
	api.Post("/refresh-token", func(c *fiber.Ctx) error {
		// Здесь должен быть реальный код для обновления токена
		return c.JSON(fiber.Map{
			"message": "Token refreshed",
		})
	})

	protected := api.Group("/protected")
	protected.Use(handlers.JWT(out.GetUsers()))

	// Маршруты, требующие аутентификации
	protected.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello, authenticated user",
		})
	})
}
