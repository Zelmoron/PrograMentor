package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"

	"main/handlers"
	"main/utils"
)

func initRoutes(app *fiber.App, in *handlers.InHandlers, out *handlers.OutHandlers) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(
			fiber.Map{
				"message": "Welcome to the API",
				"version": os.Getenv("APP_VERSION"),
			},
			os.Getenv("APP_VERSION"))
	})

	// Без аутентификации
	api := app.Group("")
	api.Post("/auth", in.Login)
	api.Post("/logout", out.LoginOut)
	api.Post("/refresh-token", out.RefreshToken)

	protected := api.Group("/protected", handlers.JWT())
	protected.Get("/", func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		tokenString := token[len("Bearer "):]
		userID, _ := utils.ValidateJWT(tokenString, os.Getenv("JWT_SECRET"))

		return c.JSON(fiber.Map{
			"message": fmt.Sprintf("User ID: %d", userID),
		})
	})
}
