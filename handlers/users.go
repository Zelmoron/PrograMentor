package handlers

import (
	"fmt"
	"net/http"
	"os"

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

	c.Locals("userID", user.ID)

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"access": token,
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
