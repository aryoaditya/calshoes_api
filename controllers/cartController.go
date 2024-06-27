package controllers

import (
	"calshoes_api/config"
	"calshoes_api/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func GetCart(c *fiber.Ctx) error {
	userClaims := c.Locals("user").(jwt.MapClaims)
	customerEmail := userClaims["email"].(string)

	// Find the customer by email
	var customer models.Customer
	if err := config.DB.Where("email = ?", customerEmail).First(&customer).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Customer not found",
			"error":   err.Error(),
		})
	}

	// Find all carts for the customer
	var carts []models.Cart
	if err := config.DB.Where("customer_id = ?", customer.Id).Preload("Customer").Find(&carts).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Cannot retrieve carts",
			"error":   err.Error(),
		})
	}

	// Success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Carts retrieved successfully",
		"data":    carts,
	})
}

func AddToCart(c *fiber.Ctx) error {
	// Get customer ID from JWT claims
	userClaims := c.Locals("user").(jwt.MapClaims)
	customerId := uint(userClaims["id"].(float64))

	// Parse request body
	type addToCartRequest struct {
		ProductId uint `json:"product_id"`
		Quantity  int  `json:"quantity"`
	}

	req := new(addToCartRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse request body",
			"error":   err.Error(),
		})
	}

	// Check if the customer already has an active cart
	var cart models.Cart
	if err := config.DB.Where("customer_id = ?", customerId).Preload("Customer").First(&cart).Error; err != nil {
		// If no active cart found, create a new one
		cart = models.Cart{
			CustomerId: customerId,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		if err := config.DB.Create(&cart).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Failed to create a new cart",
				"error":   err.Error(),
			})
		}
	} else {
		// Reload cart to get updated relations
		if err := config.DB.Preload("Customer").First(&cart, cart.Id).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Failed to load cart details",
				"error":   err.Error(),
			})
		}
	}

	// Create a new cart item
	cartItem := models.CartItem{
		CartId:    cart.Id,
		ProductId: req.ProductId,
		Quantity:  req.Quantity,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := config.DB.Create(&cartItem).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to add item to cart",
			"error":   err.Error(),
		})
	}

	// Reload cart item to get full details including relations
	if err := config.DB.Preload("Cart.Customer").Preload("Product.Category").First(&cartItem, cartItem.Id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to load cart item details",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Item added to cart successfully",
		"data":    cartItem,
	})
}
