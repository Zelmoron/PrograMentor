package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestRegistration(t *testing.T) {
	app := fiber.New()

	handler := &InHandlers{}
	app.Post("/register", handler.Registration)

	req := httptest.NewRequest("POST", "/register", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

}
