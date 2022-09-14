package v1_middleware

import (
	"os"

	"github.com/gofiber/fiber/v2"
)

func AuthCheck(c *fiber.Ctx) error {
	// check the header
	sentAuth := c.Get("Authorization")
	realAuth := os.Getenv("AUTH_TOKEN")

	if sentAuth != realAuth {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	return c.Next()
}
