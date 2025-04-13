package handlers

import (

	"fmt"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"os"
<<<<<<< HEAD

	"github.com/gofiber/fiber/v2"
=======
>>>>>>> f5937257fdca02072d4735e09ae72c074697abb3

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
<<<<<<< HEAD

	c.Locals("userID", user.ID)

=======
>>>>>>> f5937257fdca02072d4735e09ae72c074697abb3
	//TODO сделай привязку по IP к рефреш и протсо отдай его в ответе

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Logged in successfully",
	})
}

func (out *OutHandlers) LoginOut(c *fiber.Ctx) error {
	//TODO тебе не нужно здесь ничего чистить - на фронте альберт почистит localstorage и ему булет нечего отправить
	//Let do it
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

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Tokens refreshed successfully",
	})
}

func (out *OutHandlers) CheckCode(c *fiber.Ctx) error {
	userID, _ := c.Locals("userID").(int)

	var requestBody struct {
		Code string `json:"code"`
	}

	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Request parsing error",
		})
	}

	filePath, err := services.SaveUserCode(userID, requestBody.Code)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": fmt.Sprintf("Code for user ID %d saved successfully", userID)
	})
}
