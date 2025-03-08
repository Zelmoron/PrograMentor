package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

type Mock struct{}

func (in *Mock) Registration(c *fiber.Ctx) error {
	return c.SendString("Registration successful")
}

func TestRegistration(t *testing.T) {
	app := fiber.New()

	handler := &Mock{}
	app.Post("/register", handler.Registration)

	req := httptest.NewRequest("POST", "/register", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

}
