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
	insiderBuysStr := c.Query("insiderBuys")
	insiderSellsStr := c.Query("insiderSells")
	congressBuysStr := c.Query("congressBuys")
	congressSellsStr := c.Query("congressSells")

	insiderBuys := insiderBuysStr == "true"
	insiderSells := insiderSellsStr == "true"
	congressBuys := congressBuysStr == "true"
	congressSells := congressSellsStr == "true"

	orEmailInfo := bson.A{}
	if insiderBuys {
		orEmailInfo = append(orEmailInfo, bson.M{"issuerEmailInfo.insiderBuys": true})
	}
	if insiderSells {
		orEmailInfo = append(orEmailInfo, bson.M{"issuerEmailInfo.insiderSells": true})
	}
	if congressBuys {
		orEmailInfo = append(orEmailInfo, bson.M{"issuerEmailInfo.congressBuys": true})
	}
	if congressSells {
		orEmailInfo = append(orEmailInfo, bson.M{"issuerEmailInfo.congressSells": true})
	}

	if len(orEmailInfo) == 0 {
		orEmailInfo = append(orEmailInfo, bson.M{"issuerEmailInfo.noReal": true})
	}

	filter := bson.D{{Key: "issuer.issuerCik", Value: cik}, {Key: "dateAdded", Value: date},
		{Key: "$or", Value: orEmailInfo},
	}

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

	buyOrSellArray := bson.A{}
	if onBuy == "true" {
		buyOrSellArray = append(buyOrSellArray, "Buy")
	}
	if onSell == "true" {
		buyOrSellArray = append(buyOrSellArray, "Sell")
	}
	if len(buyOrSellArray) == 0 {
		buyOrSellArray = append(buyOrSellArray, "NoReal")
	}

	filter := bson.D{{Key: "reporters.reporterCik", Value: cik}, {Key: "dateAdded", Value: date},
		{Key: "buyOrSell", Value: bson.D{{Key: "$in", Value: buyOrSellArray}}},
	}

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
