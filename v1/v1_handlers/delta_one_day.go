package v1_handlers

import (
	"github.com/Matterhorn-Studios/insiderviz-forms-api/v1/v1_database"
	"github.com/Matterhorn-Studios/insidervizforms/iv_models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func DeltasFromOneDayIssuer(c *fiber.Ctx) error {
	cik := c.Query("cik")
	date := c.Query("date")

	filter := bson.D{{Key: "issuer.issuerCik", Value: cik}, {Key: "periodOfReport", Value: date}}

	cur, err := v1_database.GetCollection("DeltaForm").Find(c.Context(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err})
	}

	formList := make([]iv_models.DB_DeltaForm, 0)

	if err := cur.All(c.Context(), &formList); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err})
	}

	return c.Status(fiber.StatusOK).JSON(formList)
}

func DeltasFromOneDayReporter(c *fiber.Ctx) error {
	cik := c.Query("cik")
	date := c.Query("date")

	filter := bson.D{{Key: "reporters.reporterCik", Value: cik}, {Key: "periodOfReport", Value: date}}

	cur, err := v1_database.GetCollection("DeltaForm").Find(c.Context(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err})
	}

	formList := make([]iv_models.DB_DeltaForm, 0)

	if err := cur.All(c.Context(), &formList); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err})
	}

	return c.Status(fiber.StatusOK).JSON(formList)
}
