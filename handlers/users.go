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

	authData := AuthData{UserID: user.ID}

	token, err := utils.GenerateJWT(authData.UserID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not generate access token",
		})
	}

	refreshToken, err := utils.GenerateRefreshToken(authData.UserID)
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
		HTTPOnly: false,
		Secure:   true,
		SameSite: "Strict",
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
	})

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Logged in successfully",
	})
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
