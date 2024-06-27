package routes

import (
	"calshoes_api/controllers"
	"calshoes_api/middlewares"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	// Auth routes
	app.Post("/register", controllers.Register)
	app.Post("/login", controllers.Login)
	app.Get("/logout", controllers.Logout)

	// Grouping api endpoint
	api := app.Group("/api", middlewares.JWTMiddleware())

	//Category routes
	api.Get("/categories", controllers.GetCategories)
	api.Get("/categories/:id", controllers.GetCategoryById)

	//Products routes
	api.Post("/create_product", controllers.CreateProduct)
	api.Get("/products", controllers.GetProducts)
	api.Get("/products/:id", controllers.GetProductById)
	api.Put("/products/:id", controllers.UpdateProduct)
	api.Delete("/products/:id", controllers.DeleteProduct)

	//Products by category
	api.Get("/categories/:id/products", controllers.GetProductsByCategory)

	// Cart routes
	api.Post("/cart/add", controllers.AddToCart)
	api.Get("/cart", controllers.GetCart)
}
