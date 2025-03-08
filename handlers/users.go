package handlers

import "github.com/gofiber/fiber/v2"

func (in *InHandlers) Registration(c *fiber.Ctx) error {
	return c.SendString("Registration successful")
}

func (in *InHandlers) Login(c *fiber.Ctx) error {
	return nil
}
