package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"main/utils"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
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

	userHash := sha256.Sum256([]byte(credentials.Password))
	if hex.EncodeToString(userHash[:]) != user.Password {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid password",
		})
	}

	accessToken, err := utils.GenerateJWT(user.ID, utils.GetJWTSecret())
	if err != nil {
		return err
	}

	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		fmt.Println("Error generating refresh token:", err)
		return err
	}

	newRefreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		fmt.Println("Error generating refresh token:", err)
		return err
	}

	err = in.users.SaveRefreshToken(user.ID, refreshToken)
	if err != nil {
		fmt.Println("Error saving refresh token:", err)
		return err
	}

	return c.Status(fiber.StatusOK).JSON(LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	})
}
