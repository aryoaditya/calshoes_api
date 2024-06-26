package controllers

import (
	db "calshoes_api/config"
	"calshoes_api/models"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func CreateProduct(c *fiber.Ctx) error {
	product := new(models.Product)

	// Convert the JSON request body into the product model
	if err := c.BodyParser(product); err != nil {
		log.Printf("Error parsing product data: %v", err)
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Cannot parse JSON",
			"error":   err.Error(),
		})
	}

	// Set the time
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	// Insert the product into the database
	if err := db.DB.Create(&product).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Cannot create product",
			"error":   err.Error(),
		})
	}

	// Success CreateProduct
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data":    product,
	})
}

func GetProducts(c *fiber.Ctx) error {
	var products []models.Product

	// Get all products from the database
	if err := db.DB.Preload("Category").Find(&products).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Cannot retrieve products",
			"error":   err.Error(),
		})
	}

	// Success GetProducts
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Products retrieved successfully",
		"data":    products,
	})
}

func GetProductById(c *fiber.Ctx) error {
	// Get the product Id
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Invalid product Id",
			"error":   err.Error(),
		})
	}

	var product models.Product

	// Get the product by Id from the database
	if err := db.DB.Preload("Category").First(&product, uint(id)).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Product not found",
			"error":   err.Error(),
		})
	}

	// Success GetProductById
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Product retrieved successfully",
		"data":    product,
	})
}

func UpdateProduct(c *fiber.Ctx) error {
	// Get the product Id
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Invalid product Id",
			"error":   err.Error(),
		})
	}

	// Get the existing product from the database
	var currentProduct models.Product
	if err := db.DB.First(&currentProduct, uint(id)).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Product not found",
			"error":   err.Error(),
		})
	}

	// Convert the JSON request body into the product model
	newProduct := new(models.Product)
	if err := c.BodyParser(newProduct); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Cannot parse JSON",
			"error":   err.Error(),
		})
	}

	// Update the existing product with the new values
	currentProduct.Name = newProduct.Name
	currentProduct.Description = newProduct.Description
	currentProduct.Price = newProduct.Price
	currentProduct.ImargeUrl = newProduct.ImargeUrl
	currentProduct.CategoryId = newProduct.CategoryId
	currentProduct.UpdatedAt = time.Now()

	// Save the updated product back to the database
	if err := db.DB.Save(&currentProduct).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Cannot update product",
			"error":   err.Error(),
		})
	}

	// Success UpdateProduct
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Product updated successfully",
		"data":    currentProduct,
	})
}

func DeleteProduct(c *fiber.Ctx) error {
	// Get the product Id
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Invalid product Id",
			"error":   err.Error(),
		})
	}

	// Get the existing product from the database
	var product models.Product
	if err := db.DB.First(&product, uint(id)).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Product not found",
			"error":   err.Error(),
		})
	}

	// Delete the product from the database
	if err := db.DB.Delete(&product).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Cannot delete product",
			"error":   err.Error(),
		})
	}

	// Success DeleteProduct
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Product deleted successfully",
	})
}
