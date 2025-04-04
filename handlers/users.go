package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"

	"main/services"
	"main/utils"
)

func (in *InHandlers) Login(c *fiber.Ctx) error {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&credentials); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Request parsing error",
		})
	}

	user, err := in.repos.UsersRepo.GetUserByUsername(credentials.Username)
	if err != nil || user == nil || !services.VerifyPassword(credentials.Password, user.Password) {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid username or password",
		})
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not generate access token",
		})
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not generate refresh token",
		})
	}

	if err := in.repos.SaveRefreshToken(user.ID, refreshToken); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to save refresh token",
		})
	}
	c.Cookie(&fiber.Cookie{
		Name:     "accessToken",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 100),
		HTTPOnly: false,
		Domain:   ".ngrok-free.app",
		Secure:   false,
		SameSite: "None",
	})

	//c.Cookie(&fiber.Cookie{
	//	Name:     "refreshToken",
	//	Value:    refreshToken,
	//	Expires:  time.Now().Add(time.Hour * 100 * 100),
	//	HTTPOnly: true,
	//	Secure:   false,
	//	SameSite: "Lax",
	//})

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Logged in successfully",
	})
}

func (out *OutHandlers) LoginOut(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "accessToken",
		Value:    "",
		HTTPOnly: false,
		Secure:   false,
		SameSite: "None",
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refreshToken",
		Value:    "",
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

	updatedRefreshToken, err := utils.UpdateRefreshToken(refreshToken, userID)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "Failed to update refresh token",
		})
	}

	newAccessToken, err := utils.GenerateJWT(userID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to generate new access token",
		})
	}

	if err := out.repos.SaveRefreshToken(userID, updatedRefreshToken); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update refresh token in the database",
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "accessToken",
		Value:    newAccessToken,
		HTTPOnly: false,
		Secure:   true,
		SameSite: "Strict",
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refreshToken",
		Value:    updatedRefreshToken,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
	})

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Tokens refreshed successfully",
	})
}
