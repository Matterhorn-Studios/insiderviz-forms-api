package v1_handlers

import (
	"github.com/Matterhorn-Studios/insiderviz-forms-api/v1/models"
	"github.com/Matterhorn-Studios/insiderviz-forms-api/v1/v1_database"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func DeltasFromOneDayIssuer(c *fiber.Ctx) error {
	cik := c.Query("cik")
	date := c.Query("date")
	insiderBuys := c.Query("insiderBuys")
	insiderSells := c.Query("insiderSells")
	congressBuys := c.Query("congressBuys")
	congressSells := c.Query("congressSells")

	// TODO: Add filter for insiderBuys, insiderSells, congressBuys, congressSells

	filter := bson.D{{Key: "issuer.issuerCik", Value: cik}, {Key: "dateAdded", Value: date},
		{Key: "issuerEmailInfo.insiderBuys", Value: insiderBuys},
		{Key: "issuerEmailInfo.insiderSells", Value: insiderSells},
		{Key: "issuerEmailInfo.congressBuys", Value: congressBuys},
		{Key: "issuerEmailInfo.congressSells", Value: congressSells}}

	cur, err := v1_database.GetCollection("DeltaForm").Find(c.Context(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err})
	}

	formList := make([]models.DB_DeltaForm, 0)

	if err := cur.All(c.Context(), &formList); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err})
	}

	return c.Status(fiber.StatusOK).JSON(formList)
}

func DeltasFromOneDayReporter(c *fiber.Ctx) error {
	cik := c.Query("cik")
	date := c.Query("date")
	onBuy := c.Query("onBuy")
	onSell := c.Query("onSell")

	filter := bson.D{{Key: "reporters.reporterCik", Value: cik}, {Key: "dateAdded", Value: date}, {Key: "buyOrSell", Value: onBuy}, {Key: "buyOrSell", Value: onSell}}

	cur, err := v1_database.GetCollection("DeltaForm").Find(c.Context(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err})
	}

	formList := make([]models.DB_DeltaForm, 0)

	if err := cur.All(c.Context(), &formList); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err})
	}

	return c.Status(fiber.StatusOK).JSON(formList)
}
