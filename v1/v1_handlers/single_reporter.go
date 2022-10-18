package v1_handlers

import (
	"context"
	"math/rand"
	"sort"
	"strconv"

	"github.com/Matterhorn-Studios/insiderviz-forms-api/v1/lib"
	"github.com/Matterhorn-Studios/insiderviz-forms-api/v1/models"
	"github.com/Matterhorn-Studios/insiderviz-forms-api/v1/v1_database"
	"github.com/Matterhorn-Studios/insiderviz-forms-api/v1/v1_helpers"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func LatestThirteenF(c *fiber.Ctx) error {
	cik := c.Params("cik")

	rawOffset := c.Query("offset")
	offset := 0
	if rawOffset != "" {
		var err error
		offset, err = strconv.Atoi(rawOffset)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
		}
	}

	thirteenFCollection := v1_database.GetCollection("13F")

	filter := bson.D{{Key: "cik", Value: cik}}
	opts := options.FindOne().SetSort(bson.D{{Key: "periodOfReport", Value: -1}}).SetSkip(int64(offset))

	form := thirteenFCollection.FindOne(context.TODO(), filter, opts)

	var thirteenF models.DB_Form13F_Base
	if form.Err() != nil {
		return c.JSON(fiber.Map{"status": "empty", "form": models.DB_Form13F_Base{}})
	} else {
		err := form.Decode(&thirteenF)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
		} else {
			// update the document so that only 40 companies show up
			if len(thirteenF.Holdings) > 40 {
				// sort the holdings by value descending
				sort.Slice(thirteenF.Holdings, func(i, j int) bool {
					return thirteenF.Holdings[i].NetTotal > thirteenF.Holdings[j].NetTotal
				})

				otherTotal := 0.0
				otherShares := 0.0
				// remove the last holdings
				for i := 40; i < len(thirteenF.Holdings); i++ {
					otherTotal += float64(thirteenF.Holdings[i].NetTotal)
					otherShares += float64(thirteenF.Holdings[i].Shares)
				}

				thirteenF.Holdings = thirteenF.Holdings[:40]
				thirteenF.Holdings = append(thirteenF.Holdings, models.DB_Form13F_Holding{
					Name:     "Other",
					NetTotal: float32(otherTotal),
					Shares:   float32(otherShares),
					Cik:      "",
				})
			}

			return c.JSON(fiber.Map{"status": "ok", "form": thirteenF})
		}
	}

}

func Reporter(c *fiber.Ctx) error {
	cik := c.Params("cik")

	filter := bson.D{{Key: "reporters.reporterCik", Value: cik}, {
		Key: "$or", Value: bson.A{
			bson.D{{Key: "formClass", Value: "Insider"}},
			bson.D{{Key: "formClass", Value: "Congress"}},
		},
	}}
	opts := options.Find().SetSort(bson.D{{Key: "periodOfReport", Value: -1}})

	deltaForms, err := lib.DeltaFormFetch(filter, opts)

	if err != nil {
		deltaForms = []models.DB_DeltaForm{}
	}

	// get the reporter's information
	reporterCollection := v1_database.GetCollection("Reporter")
	filter = bson.D{{Key: "cik", Value: cik}}

	issuerInfo := reporterCollection.FindOne(context.TODO(), filter)

	if issuerInfo.Err() != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": issuerInfo.Err()})
	}

	// unmarshal the reporter info
	var reporter bson.M
	if err = issuerInfo.Decode(&reporter); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}

	for i := 0; i < len(deltaForms); i++ {
		deltaForms[i].PercentChange = v1_helpers.PercentChange(deltaForms[i])
	}

	return c.JSON(fiber.Map{"forms": deltaForms, "info": reporter})
}

func RandomReporter(c *fiber.Ctx) error {
	// get the reporter collection
	reporterCollection := v1_database.GetCollection("Reporter")
	var reporter models.DB_Reporter_Doc

	for {
		// get a random offset 1-5000
		offset := rand.Intn(5000) + 1
		opts := options.FindOne().SetSkip(int64(offset))

		cursor := reporterCollection.FindOne(context.TODO(), bson.D{}, opts)
		if cursor.Err() == nil {
			// decode into DB_Issuer_Doc
			err := cursor.Decode(&reporter)
			if err == nil {
				// check the reporter has forms
				filter := bson.D{{Key: "reporters.reporterCik", Value: reporter.Cik}}
				opts2 := options.Find().SetSort(bson.D{{Key: "periodOfReport", Value: -1}}).SetLimit(1)
				deltaForms, err := lib.DeltaFormFetch(filter, opts2)
				if err == nil {
					if len(deltaForms) > 0 && deltaForms[0].PeriodOfReport > "2021-01-01" {
						return c.JSON(reporter)
					}
				}
			}
		}

	}
}
