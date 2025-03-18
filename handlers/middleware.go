package handlers

import (
	"main/services"
	"main/utils"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func authMiddleware(secretKey string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Missing authorization header",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid authorization header format",
			})
		}

		tokenString := parts[1]

		sub, err := utils.ValidateJWT(tokenString, secretKey)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid token",
			})
		}

		c.Locals("sub", sub)
		return c.Next()
	}
}

func JWT(users *services.Users) fiber.Handler {
	return authMiddleware(os.Getenv("JWT_SECRET"))
}

func REFRESH(users *services.Users) fiber.Handler {
	return authMiddleware(os.Getenv("REFRESH_SECRET"))
}
