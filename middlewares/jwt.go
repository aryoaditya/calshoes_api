package middlewares

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// JWTMiddleware is the middleware to protect routes
func JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the token from cookie
		tokenString := c.Cookies("token")
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "You are not authorized",
			})
		}

		// Parse and validate the token
		claims := &jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "Unexpected signing method")
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			v, _ := err.(*jwt.ValidationError)
			switch v.Errors {
			case jwt.ValidationErrorSignatureInvalid:
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"success": false,
					"message": "You are not authorized",
				})
			case jwt.ValidationErrorExpired:
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"success": false,
					"message": "You are not authorized, Token expired!",
				})
			default:
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"success": false,
					"message": "You are not authorized",
				})
			}
		}

		if !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "You are not authorized",
			})
		}

		c.Locals("user", *claims)
		return c.Next()
	}
}

// GenerateJWT generates a new JWT token
func GenerateJWT(email string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(), // Token expires in 72 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
