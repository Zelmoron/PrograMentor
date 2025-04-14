package handlers

import (
	"context"
	"net/http"
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

	//TODO Переход к сервису, который будет создавать докер
	logs := make(chan string)
	errChan := make(chan error)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	go services.StartUserCode(ctx, logs, errChan, filePath)

	select {
	case log := <-logs:
		return c.Status(http.StatusOK).JSON(fiber.Map{
			"message": log,
		})
	case err := <-errChan:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	case <-ctx.Done():
		if ctx.Err() == context.DeadlineExceeded {
			return c.Status(http.StatusRequestTimeout).JSON(fiber.Map{
				"message": "Error: Waiting time finally exceeded",
			})
		}
	}

	return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
		"message": "Unknown error",
	})
}
