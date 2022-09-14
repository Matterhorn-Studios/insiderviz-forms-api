package v1_handlers

import (
	"context"
	"math/rand"
	"sync"

	"github.com/Matterhorn-Studios/insiderviz-forms-api/v1/lib"
	"github.com/Matterhorn-Studios/insiderviz-forms-api/v1/v1_database"
	"github.com/Matterhorn-Studios/insidervizforms/iv_models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Issuer(c *fiber.Ctx) error {
	cik := c.Params("cik")
	includeGraph := c.Query("includeGraph")

	filter := bson.D{{Key: "issuer.issuerCik", Value: cik}, {
		Key: "$or", Value: bson.A{
			bson.D{{Key: "formClass", Value: "Insider"}},
			bson.D{{Key: "formClass", Value: "Congress"}},
		},
	}}
	opts := options.Find().SetSort(bson.D{{Key: "periodOfReport", Value: -1}})

	var wg sync.WaitGroup
	var deltaForms []iv_models.DB_DeltaForm
	wg.Add(2)
	go func() {
		defer wg.Done()
		var err error
		deltaForms, err = lib.DeltaFormFetch(filter, opts)
		if err != nil {
			deltaForms = []iv_models.DB_DeltaForm{}
		}
	}()

	var issuer iv_models.DB_Issuer_Doc
	var err error
	go func() {
		defer wg.Done()
		if includeGraph == "true" {
			issuer, err = lib.UpdateStockData(cik)

		} else {
			// get the issuer's information
			issuerCollection := v1_database.GetCollection("Issuer")
			filter = bson.D{{Key: "cik", Value: cik}}
			opts := options.FindOne().SetProjection(bson.D{{Key: "stockData", Value: 0}})

			// do not include the stockData
			issuerInfo := issuerCollection.FindOne(context.TODO(), filter, opts)

			if issuerInfo.Err() != nil {
				err = issuerInfo.Err()
			}

			// unmarshal the issuer info
			err = issuerInfo.Decode(&issuer)
		}
	}()

	wg.Wait()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}

	return c.JSON(fiber.Map{"forms": deltaForms, "info": issuer})
}

func RandomIssuer(c *fiber.Ctx) error {
	// get the issuer collection
	issuerCollection := v1_database.GetCollection("Issuer")
	var issuer iv_models.DB_Issuer_Doc

	for {
		// get a random offset 1-5000
		offset := rand.Intn(5000) + 1
		opts := options.FindOne().SetSkip(int64(offset)).SetProjection(bson.D{{Key: "stockData", Value: 0}})

		cursor := issuerCollection.FindOne(context.TODO(), bson.D{}, opts)
		if cursor.Err() == nil {
			// decode into DB_Issuer_Doc
			err := cursor.Decode(&issuer)
			if err == nil {
				if len(issuer.Tickers) > 0 {
					// check the issuer has forms
					filter := bson.D{{Key: "issuer.issuerCik", Value: issuer.Cik}}
					opts2 := options.Find().SetSort(bson.D{{Key: "periodOfReport", Value: -1}}).SetLimit(1)
					deltaForms, err := lib.DeltaFormFetch(filter, opts2)
					if err == nil {
						if len(deltaForms) > 0 && deltaForms[0].PeriodOfReport > "2021-01-01" {
							return c.JSON(issuer)
						}
					}
				}
			}
		}
	}

}

func IssuerGraph(c *fiber.Ctx) error {
	cik := c.Params("cik")

	// get the stock data
	stockData, err := lib.UpdateStockData(cik)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}

	return c.JSON(stockData)
}
