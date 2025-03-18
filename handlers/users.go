package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"main/utils"
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

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		return err
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:     "accessToken",
		Value:    token,
		HTTPOnly: false,
		Secure:   false,
		SameSite: "None",
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		HTTPOnly: true,
		Secure:   false,
		SameSite: "None",
	})

	c.Status(http.StatusOK)
	return nil
}
