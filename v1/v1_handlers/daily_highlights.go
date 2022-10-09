package v1_handlers

import (
	"github.com/Matterhorn-Studios/insiderviz-forms-api/v1/models"
	"github.com/Matterhorn-Studios/insiderviz-forms-api/v1/v1_database"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type CompanySum struct {
	Cik    string
	Ticker string
	Total  float32
}

func DailyHighlights(c *fiber.Ctx) error {
	date := c.Query("date")

	// get all of the trades added today
	filter := bson.D{{Key: "dateAdded", Value: date}}

	cur, err := v1_database.GetCollection("DeltaForm").Find(c.Context(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err})
	}

	formList := make([]models.DB_DeltaForm, 0)
	if err = cur.All(c.Context(), &formList); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err})
	}

	numTrades := len(formList)
	largestBuy := float32(0.0)
	largestSell := float32(0.0)
	highestVolume := ""
	companySums := make([]CompanySum, 0)

	for _, form := range formList {
		if form.NetTotal > largestBuy && form.BuyOrSell == "Buy" {
			largestBuy = form.NetTotal
		}
		if form.NetTotal > largestSell && form.BuyOrSell == "Sell" {
			largestSell = form.NetTotal
		}

		// check if the company is already in the list
		found := false
		for i, company := range companySums {
			if company.Cik == form.Issuer.IssuerCik {
				companySums[i].Total += form.NetTotal
				found = true
				break
			}
		}
		if !found {
			companySums = append(companySums, CompanySum{Cik: form.Issuer.IssuerCik, Ticker: form.Issuer.IssuerTicker, Total: form.NetTotal})
		}
	}

	// get the company with the highest volume
	if len(companySums) > 0 {

		highestVolume = companySums[0].Ticker
		for _, company := range companySums {
			if company.Total > companySums[0].Total {
				highestVolume = company.Ticker
			}
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"numTrades":     numTrades,
		"largestBuy":    largestBuy,
		"largestSell":   largestSell,
		"highestVolume": highestVolume,
	})
}
