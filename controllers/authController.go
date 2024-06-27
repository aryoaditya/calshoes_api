package controllers

import (
	"calshoes_api/config"
	"calshoes_api/middlewares"
	"calshoes_api/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// Register a new user
func Register(c *fiber.Ctx) error {
	user := new(models.Customer)

	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Cannot parse JSON",
			"error":   err.Error(),
		})
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Cannot hash password",
			"error":   err.Error(),
		})
	}

	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Insert to database
	if err := config.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Cannot create account",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Account registered successfully",
	})
}

// Login user
func Login(c *fiber.Ctx) error {
	var data map[string]string

	// Input JSON req
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Cannot parse JSON",
			"error":   err.Error(),
		})
	}

	var user models.Customer
	if err := config.DB.Where("email = ?", data["email"]).First(&user).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "User not found",
			"error":   err.Error(),
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data["password"])); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid email or password",
			"error":   err.Error(),
		})
	}

	token, err := middlewares.GenerateJWT(user.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Could not login",
			"error":   err.Error(),
		})
	}

	// set token to cookie
	cookie := fiber.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(72 * time.Hour), // Token expires in 72 hours
		HTTPOnly: true,
		SameSite: "None",
	}

	c.Cookie(&cookie)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Login successful",
		"token":   token,
	})
}

// Logout handles user logout and clears JWT token from cookie
func Logout(c *fiber.Ctx) error {
	// Clear the token cookie
	cookie := fiber.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour), // Set expiration in the past to delete cookie
		HTTPOnly: true,
		SameSite: "None",
	}

	c.Cookie(&cookie)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Logout successful",
	})
}
