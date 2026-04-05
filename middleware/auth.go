package middleware

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{"error": "Unauthorized, butuh token!"})
		}

		tokenString := authHeader[7:] // Menghapus "Bearer "
		token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if token != nil && token.Valid {
			return c.Next()
		}
		return c.Status(401).JSON(fiber.Map{"error": "Token tidak valid!"})
	}
}