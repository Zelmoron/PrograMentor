package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (in *InHandlers) Login(c *fiber.Ctx) error {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&credentials); err != nil {
		return err
	}

	user, err := in.repos.UsersRepo.GetUserByUsername(credentials.Username)
	if err != nil {
		return err
	}

	userHash := sha256.Sum256([]byte(credentials.Password))
	if hex.EncodeToString(userHash[:]) != user.Password {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid password",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Login successful!",
	})
}
