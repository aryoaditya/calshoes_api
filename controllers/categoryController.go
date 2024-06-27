package controllers

import (
	"calshoes_api/config"
	"calshoes_api/models"

	"github.com/gofiber/fiber/v2"
)

func GetCategories(c *fiber.Ctx) error {
	var categories []models.Category

	// Get all categories from the database
	if err := config.DB.Preload("Products").Find(&categories).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Cannot retrieve categories",
			"error":   err.Error(),
		})
	}

	// Success GetCategories
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Categories retrieved successfully",
		"data":    categories,
	})
}

func GetCategoryById(c *fiber.Ctx) error {
	// Get category Id
	idParam := c.Params("id")
	var category models.Category

	// Find category by Id
	if err := config.DB.Preload("Products").First(&category, idParam).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Category not found",
			"error":   err.Error(),
		})
	}

	// Success GetCategoryById
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Category retrieved successfully",
		"data":    category,
	})
}
