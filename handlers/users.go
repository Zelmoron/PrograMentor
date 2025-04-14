package handlers

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
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

	err := services.SaveUserCode(userID, requestBody.Code)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	//TODO Переход к сервису, который будет создавать докер
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.41"))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	filename := fmt.Sprintf("./codes/%d.go", userID)
	defer os.Remove(filename)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "golang:latest",
		Cmd:   []string{"go", "run", "/code/code.go"},
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: filename,
				Target: "/code/code.go",
			},
		},
	}, nil, nil, "")
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	statusCh, erCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)

	select {
	case err := <-erCh:
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
	case <-statusCh:
	}

	logs, err := cli.ContainerLogs(ctx, resp.ID, container.LogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	var buf bytes.Buffer
	buf.ReadFrom(logs)

	if err := cli.ContainerStop(ctx, resp.ID, container.StopOptions{}); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to stop container",
		})
	}
	if err := cli.ContainerRemove(ctx, resp.ID, container.RemoveOptions{}); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to remove container",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": buf.String(),
	})
}
