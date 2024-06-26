package main

import (
	db "calshoes_api/config"
	"calshoes_api/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	db.Connect()
	// Create an instance
	app := fiber.New()

	// Define a route
	routes.Setup(app)

	// Start the Fiber app
	app.Listen(":3000")
}
