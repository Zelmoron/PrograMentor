package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"main/services"
	"main/utils"
)

type AuuthData struct {
	UserID int64
}

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

	if !services.VerifyPassword(credentials.Password, user.Password) {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid password",
		})
	}

	c.Locals("authData", AuuthData{
		UserID: user.ID,
	})

	return c.Next()
}

func (out *OutHandlers) LoginOut(c *fiber.Ctx) error {
	authData, _ := c.Locals("authData").(AuuthData)

	token, err := utils.GenerateJWT(authData.UserID)
	if err != nil {
		return err
	}

	refreshToken, err := utils.GenerateRefreshToken(authData.UserID)
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

	return c.SendStatus(fiber.StatusOK)
}
