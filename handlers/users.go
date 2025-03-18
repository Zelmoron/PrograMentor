package handlers

import (
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"

	"main/services"
	"main/utils"
)

type AuthData struct {
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

	c.Locals("authData", AuthData{
		UserID: user.ID,
	})

	return c.Next()
}

func (out *OutHandlers) LoginOut(c *fiber.Ctx) error {
	authData, _ := c.Locals("authData").(AuthData)

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

func (out *OutHandlers) RefreshToken(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refreshToken")

	if refreshToken == "" {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "Missing refresh token",
		})
	}

	userID, err := utils.ValidateRefreshToken(refreshToken, os.Getenv("REFRESH_SECRET"))
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid refresh token",
		})
	}

	newAccessToken, err := utils.GenerateJWT(userID)
	if err != nil {
		return err
	}
	newRefreshToken, err := utils.GenerateRefreshToken(userID)
	if err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:     "accessToken",
		Value:    newAccessToken,
		HTTPOnly: false,
		Secure:   false,
		SameSite: "None",
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refreshToken",
		Value:    newRefreshToken,
		HTTPOnly: true,
		Secure:   false,
		SameSite: "None",
	})

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Tokens refreshed successfully",
	})
}
