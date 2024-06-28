package main

import (
	"calshoes_api/config"
	"calshoes_api/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config.ConnectDB()

	// Create an instance
	app := fiber.New()

	// Define a route
	routes.Setup(app)

	// Start the Fiber app
	app.Listen(":8080")
}
