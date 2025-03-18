package main

import (
	"os"

	"github.com/gofiber/fiber/v2"

	"main/handlers"
	"main/services"
)

func initRoutes(app *fiber.App, in *handlers.InHandlers, out *handlers.OutHandlers, users *services.Users) {
	app.Get("/", func(c *fiber.Ctx) error { return c.JSON(os.Getenv("APP_VERSION")) })

	// Без аутентификации
	api := app.Group("")
	api.Post("/auth", in.Login, out.LoginOut)
	api.Post("/refresh-token", in.RefreshToken, out.RefreshTokenOut)

	protected := api.Group("/protected")
	protected.Use(handlers.JWT(users))

	// Маршруты, требующие аутентификации
	protected.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello, authenticated user",
		})
	})
}
