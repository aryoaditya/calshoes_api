package routes

import (
	"calshoes_api/controllers"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	api := app.Group("/api")

	//Category routes
	api.Get("/categories", controllers.GetCategories)
	api.Get("/categories/:id", controllers.GetCategoryById)

	//Products routes
	api.Post("/create_product", controllers.CreateProduct)
	api.Get("/products", controllers.GetProducts)
	api.Get("/products/:id", controllers.GetProductById)
	api.Put("/products/:id", controllers.UpdateProduct)
	api.Delete("/products/:id", controllers.DeleteProduct)
}
