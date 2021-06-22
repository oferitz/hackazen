package main

import (
	"github.com/gofiber/fiber/v2"
)

func (app *application) healthHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":      "available",
		"environment": app.config.String("app.environment"),
		"version":     app.config.String("app.version"),
	})

}
