package controllers

import (
	"calshoes_api/config"
	"calshoes_api/models"
	"time"

	"github.com/gofiber/fiber/v2"
)

func ProcessPayment(c *fiber.Ctx) error {
	var req models.Payment
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse request body",
			"error":   err.Error(),
		})
	}

	// Create a new payment record
	newPayment := models.Payment{
		OrderId:   req.OrderId,
		Method:    req.Method,
		Amount:    req.Amount,
		Status:    req.Status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tx := config.DB.Begin()

	if err := tx.Create(&newPayment).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to process payment",
			"error":   err.Error(),
		})
	}

	// Update order status
	if err := tx.Model(&models.Order{}).Where("id = ?", newPayment.OrderId).Update("status", newPayment.Status).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to update order status",
			"error":   err.Error(),
		})
	}

	// Commit transaction
	tx.Commit()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Payment processed successfully",
		"data":    newPayment,
	})
}
