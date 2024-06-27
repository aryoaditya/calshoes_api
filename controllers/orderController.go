package controllers

import (
	"calshoes_api/config"
	"calshoes_api/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func GetOrders(c *fiber.Ctx) error {
	userClaims := c.Locals("user").(jwt.MapClaims)
	customerId := uint(userClaims["id"].(float64))

	var orders []models.Order
	if err := config.DB.Preload("Customer").Where("customer_id = ?", customerId).Find(&orders).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to retrieve orders",
			"error":   err.Error(),
		})
	}

	// Success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Orders retrieved successfully",
		"data":    orders,
	})
}
