package tests

import (
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

type InHandlers struct{}

func (in *InHandlers) Registration(c *fiber.Ctx) error {
	return nil
}

func TestRegistration(t *testing.T) {

	app := fiber.New()

	handler := &InHandlers{}
	app.Post("/register", handler.Registration)

	req := httptest.NewRequest("POST", "/register", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}
