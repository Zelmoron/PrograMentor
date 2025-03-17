package handlers

import (
	"main/services"
	"main/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(users *services.Users) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("auth")
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

		userID, err := ValidateJWTFromService(tokenString, users)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid token",
			})
		}

		c.Locals("userID", userID)
		return c.Next()
	}
}

func ValidateJWTFromService(tokenString string, users *services.Users) (uint, error) {
	return utils.ValidateJWT(tokenString, users.Cfg.JWTSecret)
}
