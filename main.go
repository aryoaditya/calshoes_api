package main

import (
	db "calshoes_api/config"

	"github.com/gofiber/fiber/v2"
)

func main() {
	db.Connect()
	// Create an instance
	app := fiber.New()

	// Define a route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello, World!",
		})
	})

	// Start the Fiber app
	app.Listen(":3000")
}
