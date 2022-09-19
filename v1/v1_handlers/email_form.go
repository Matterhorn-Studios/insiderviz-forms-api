package v1_handlers

import (
	"github.com/Matterhorn-Studios/insiderviz-forms-api/v1/lib"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func EmailFormIssuer(c *fiber.Ctx) error {
	cik := c.Params("cik")
	date := c.Params("date")

	filter := bson.D{{Key: "issuer.issuerCik", Value: cik}, {Key: "periodOfReport", Value: date}}
	opts := options.Find()

	forms, err := lib.DeltaFormFetch(filter, opts)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}

	return c.JSON(forms)
}

func EmailFormReporter(c *fiber.Ctx) error {
	cik := c.Params("cik")
	date := c.Params("date")

	filter := bson.D{{Key: "reporters.reporterCik", Value: bson.D{{Key: "$all", Value: bson.A{cik}}}}, {Key: "periodOfReport", Value: date}}
	opts := options.Find()

	forms, err := lib.DeltaFormFetch(filter, opts)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}

	return c.JSON(forms)
}
