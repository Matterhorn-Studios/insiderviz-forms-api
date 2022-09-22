package v1_handlers

import "github.com/gofiber/fiber/v2"

func DeltaFromCsv(c *fiber.Ctx) error {
	cik := c.Params("cik")
	return c.JSON(fiber.Map{"test": cik})
}
