package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header is required",
			})
		}

		// Check if it's a Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization format must be Bearer {token}",
			})
		}

		token := parts[1]

		// Validate the token (implement token validation logic)
		// For now, just set dummy values - in production, validate the JWT token
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		// to do: jwt validation
		// For now, just set dummy values
		userID, _ := primitive.ObjectIDFromHex("000000000000000000000000")
		c.Locals("userID", userID)
		c.Locals("username", "testuser")
		c.Locals("role", "user")

		return c.Next()
	}
}
