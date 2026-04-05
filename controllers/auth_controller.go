package controllers

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func Login(c *fiber.Ctx) error {
	type LoginInput struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var input LoginInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Input salah"})
	}

	// Cek kredensial dari .env
	if input.Username != os.Getenv("ADMIN_USER") || input.Password != os.Getenv("ADMIN_PASS") {
		return c.Status(401).JSON(fiber.Map{"error": "Username atau Password salah"})
	}

	// Buat JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": input.Username,
		"exp":  time.Now().Add(time.Hour * 72).Unix(),
	})

	t, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	return c.JSON(fiber.Map{"token": t})
}