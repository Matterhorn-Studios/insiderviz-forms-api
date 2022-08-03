package delta

import (
	"net/http"
	"strconv"

	"github.com/Matterhorn-Studios/insiderviz-forms-api/v1/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Recent(c *gin.Context) {
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

	deltaForms, err := utils.DeltaFormFetch(filter, opts)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error1": err.Error()})
		return
	}

	c.JSON(http.StatusOK, deltaForms)

}
