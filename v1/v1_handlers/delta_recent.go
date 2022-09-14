package v1_handlers

import (
	"strconv"

	"github.com/Matterhorn-Studios/insiderviz-forms-api/v1/lib"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Recent(c *fiber.Ctx) error {
	// get the params
	buySellOrBoth := c.Query("buySellOrBoth")
	insiderCongressOrBoth := c.Query("insiderCongressOrBoth")

	var filterOne bson.E = bson.E{}
	if buySellOrBoth == "Buy" {
		filterOne = bson.E{Key: "buyOrSell", Value: "Buy"}
	} else if buySellOrBoth == "Sell" {
		filterOne = bson.E{Key: "buyOrSell", Value: "Sell"}
	}

	var filterTwo bson.E = bson.E{}
	if insiderCongressOrBoth == "Insider" {
		filterTwo = bson.E{Key: "formClass", Value: "Insider"}
	} else if insiderCongressOrBoth == "Congress" {
		filterTwo = bson.E{Key: "formClass", Value: "Congress"}
	}

	// offset and limit
	offset := c.Query("offset")
	limit := c.Query("limit")

	// parse the offset and limit
	var offsetInt int64 = 0
	var limitInt int64 = 50

	if offset != "" {
		offsetInt, _ = strconv.ParseInt(offset, 10, 64)
	}

	if limit != "" {
		limitInt, _ = strconv.ParseInt(limit, 10, 64)
	}

	filter := bson.D{filterOne, filterTwo}
	opts := options.Find().SetLimit(limitInt).SetSkip(offsetInt).SetSort(bson.D{{Key: "periodOfReport", Value: -1}})

	deltaForms, err := lib.DeltaFormFetch(filter, opts)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}

	return c.JSON(deltaForms)

}
