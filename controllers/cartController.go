package controllers

import (
	"calshoes_api/config"
	"calshoes_api/models"
	"strconv"
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

	// Check if the product already exists
	var existingCartItem models.CartItem
	if err := config.DB.Where("cart_id = ? AND product_id = ?", cart.Id, req.ProductId).First(&existingCartItem).Error; err != nil {
		// If not found, create a new cart item
		newCartItem := models.CartItem{
			CartId:    cart.Id,
			ProductId: req.ProductId,
			Quantity:  req.Quantity,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := config.DB.Create(&newCartItem).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Failed to add item to cart",
				"error":   err.Error(),
			})
		}

		// Reload cart item to get full details
		if err := config.DB.Preload("Cart.Customer").Preload("Product.Category").First(&newCartItem, newCartItem.Id).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Failed to load cart item details",
				"error":   err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "Item added to cart successfully",
			"data":    newCartItem,
		})
	}

	// If found, update the quantity of existing cart item
	existingCartItem.Quantity += req.Quantity
	existingCartItem.UpdatedAt = time.Now()

	if err := config.DB.Save(&existingCartItem).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to update cart item quantity",
			"error":   err.Error(),
		})
	}

	// Reload cart item to get full details
	if err := config.DB.Preload("Cart.Customer").Preload("Product.Category").First(&existingCartItem, existingCartItem.Id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to load cart item details",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Item quantity updated in cart successfully",
		"data":    existingCartItem,
	})
}

func GetCartItems(c *fiber.Ctx) error {
	userClaims := c.Locals("user").(jwt.MapClaims)
	customerId := uint(userClaims["id"].(float64))

	// Find the customer's cart
	var cart models.Cart
	if err := config.DB.Where("customer_id = ?", customerId).Preload("Customer").First(&cart).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Cart not found",
			"error":   err.Error(),
		})
	}

	// Find all cart items
	var cartItems []models.CartItem
	if err := config.DB.Where("cart_id = ?", cart.Id).Preload("Cart.Customer").Preload("Product.Category").Find(&cartItems).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to retrieve cart items",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Cart items retrieved successfully",
		"data":    cartItems,
	})
}

func DeleteCartItem(c *fiber.Ctx) error {
	userClaims := c.Locals("user").(jwt.MapClaims)
	customerId := uint(userClaims["id"].(float64))

	// Parse product ID
	productID, err := strconv.ParseUint(c.Params("product_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid product ID",
			"error":   err.Error(),
		})
	}

	// Check if the customer has an active cart
	var cart models.Cart
	if err := config.DB.Where("customer_id = ?", customerId).First(&cart).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Cart not found",
			"error":   err.Error(),
		})
	}

	// Check if the cart item exists
	var cartItem models.CartItem
	if err := config.DB.Where("cart_id = ? AND product_id = ?", cart.Id, productID).First(&cartItem).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Cart item not found",
			"error":   err.Error(),
		})
	}

	// Delete the cart item from the database
	if err := config.DB.Delete(&cartItem).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to delete cart item",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Cart item deleted successfully",
	})
}

func Checkout(c *fiber.Ctx) error {
	userClaims := c.Locals("user").(jwt.MapClaims)
	customerId := uint(userClaims["id"].(float64))

	var cart models.Cart
	if err := config.DB.Where("customer_id = ?", customerId).Preload("Customer").First(&cart).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Cart not found",
			"error":   err.Error(),
		})
	}

	var cartItems []models.CartItem
	if err := config.DB.Where("cart_id = ?", cart.Id).Preload("Product").Find(&cartItems).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to retrieve cart items",
			"error":   err.Error(),
		})
	}

	var totalPrice float64
	for _, item := range cartItems {
		totalPrice += float64(item.Quantity) * item.Product.Price
	}

	newOrder := models.Order{
		CustomerId: customerId,
		TotalPrice: totalPrice,
		Status:     "Pending",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := config.DB.Create(&newOrder).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to create a new order",
			"error":   err.Error(),
		})
	}

	if err := config.DB.Where("cart_id = ?", cart.Id).Delete(&models.CartItem{}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to clear cart items",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Checkout successful",
		"data":    newOrder,
	})
}
