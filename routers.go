package main

import (
	"os"

	"github.com/gofiber/fiber/v2"

	"main/handlers"
)

func initRoutes(app *fiber.App, in *handlers.InHandlers, out *handlers.OutHandlers) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(
			os.Getenv("APP_VERSION"))
	})

	// Без аутентификации
	api := app.Group("")
	api.Post("/auth", in.Login)

	protected := api.Group("/protected", handlers.JWT())

	protected.Post("/check-code", out.CheckCode)
	//TODO добавить пагинацию
}
