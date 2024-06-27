package controllers

import (
	"calshoes_api/config"
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Cannot parse JSON",
			"error":   err.Error(),
		})
	}

	// Set the time
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	// Insert the product into the database
	if err := config.DB.Create(&product).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Cannot create product",
			"error":   err.Error(),
		})
	}

	// Success CreateProduct
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data":    product,
	})
}

func GetProducts(c *fiber.Ctx) error {
	var products []models.Product

	// Get all products from the database
	if err := config.DB.Preload("Category").Find(&products).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Cannot retrieve products",
			"error":   err.Error(),
		})
	}

	// Success GetProducts
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid product Id",
			"error":   err.Error(),
		})
	}

	var product models.Product

	// Get the product by Id from the database
	if err := config.DB.Preload("Category").First(&product, uint(id)).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Product not found",
			"error":   err.Error(),
		})
	}

	// Success GetProductById
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid product Id",
			"error":   err.Error(),
		})
	}

	// Get the existing product from the database
	var currentProduct models.Product
	if err := config.DB.First(&currentProduct, uint(id)).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Product not found",
			"error":   err.Error(),
		})
	}

	// Convert the JSON request body into the product model
	newProduct := new(models.Product)
	if err := c.BodyParser(newProduct); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Cannot parse JSON",
			"error":   err.Error(),
		})
	}

	// Update the existing product with the new values
	currentProduct.Name = newProduct.Name
	currentProduct.Description = newProduct.Description
	currentProduct.Price = newProduct.Price
	currentProduct.ImageUrl = newProduct.ImageUrl
	currentProduct.CategoryId = newProduct.CategoryId
	currentProduct.UpdatedAt = time.Now()

	// Save the updated product back to the database
	if err := config.DB.Save(&currentProduct).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Cannot update product",
			"error":   err.Error(),
		})
	}

	// Success UpdateProduct
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid product Id",
			"error":   err.Error(),
		})
	}

	// Get the existing product from the database
	var product models.Product
	if err := config.DB.First(&product, uint(id)).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Product not found",
			"error":   err.Error(),
		})
	}

	// Delete the product from the database
	if err := config.DB.Delete(&product).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Cannot delete product",
			"error":   err.Error(),
		})
	}

	// Success DeleteProduct
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Product deleted successfully",
	})
}

func GetProductsByCategory(c *fiber.Ctx) error {
	categoryId := c.Params("id")
	var products []models.Product
	var category models.Category

	// Check if the category exists
	if err := config.DB.First(&category, categoryId).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Category not found",
		})
	}

	// Find products by category Id
	if err := config.DB.Where("category_id = ?", categoryId).Find(&products).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Cannot retrieve products",
			"error":   err.Error(),
		})
	}

	// If no products found for the category_id
	if len(products) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "No products found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Products retrieved successfully",
		"data":    products,
	})
}
