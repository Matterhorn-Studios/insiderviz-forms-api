package delta

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Matterhorn-Studios/insiderviz-forms-api/v1/structs"
	"github.com/Matterhorn-Studios/insiderviz-forms-api/v1/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// LIST OF PARAMS
// sortBy: periodOfReport, netTotal, shares, sharePrice
// formClass: Insider, Congress, Institution
// buyOrSellOrBoth: Buy, Sell, Both
// order: asc, desc
// dateStart: YYYY-MM-DD
// dateEnd: YYYY-MM-DD
// netTotalMin: float64
// netTotalMax: float64
// sharesMin: float64
// sharesMax: float64
// sharePriceMin: float64
// sharePriceMax: float64
// take: int
// skip: int

func DeepFilter(c *gin.Context) {
	// get the params and validate
	sortBy := c.Query("sortBy")
	if sortBy != "periodOfReport" && sortBy != "netTotal" && sortBy != "shares" && sortBy != "sharePrice" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "sortBy must be periodOfReport, netTotal, shares, or sharePrice"})
		return
	}

	formClass := c.Query("formClass")
	formClassArray := strings.Split(formClass, ",")

	order := c.Query("order")
	if order != "asc" && order != "desc" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order must be asc or desc"})
	}

	dateStart := c.Query("dateStart")
	dateEnd := c.Query("dateEnd")

	buyOrSellOrBoth := c.Query("buyOrSellOrBoth")

	netTotalMinQuery := c.Query("netTotalMin")
	netTotalMin, err := strconv.ParseFloat(netTotalMinQuery, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "netTotalMin must be a float"})
		return
	}

	netTotalMaxQuery := c.Query("netTotalMax")
	netTotalMax, err := strconv.ParseFloat(netTotalMaxQuery, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "netTotalMax must be a float"})
		return
	}

	sharesMinQuery := c.Query("sharesMin")
	sharesMin, err := strconv.ParseFloat(sharesMinQuery, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "sharesMin must be a float"})
		return
	}

	sharesMaxQuery := c.Query("sharesMax")
	sharesMax, err := strconv.ParseFloat(sharesMaxQuery, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "sharesMax must be a float"})
		return
	}

	sharePriceMinQuery := c.Query("sharePriceMin")
	sharePriceMin, err := strconv.ParseFloat(sharePriceMinQuery, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "sharePriceMin must be a float"})
		return
	}

	sharePriceMaxQuery := c.Query("sharePriceMax")
	sharePriceMax, err := strconv.ParseFloat(sharePriceMaxQuery, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "sharePriceMax must be a float"})
		return
	}

	takeQuery := c.Query("take")
	take, err := strconv.ParseInt(takeQuery, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "take must be an int"})
		return
	}

	skipQuery := c.Query("skip")
	skip, err := strconv.ParseInt(skipQuery, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "skip must be an int"})
		return
	}

	// call the function
	forms, err := HandleDeepFilter(sortBy, order, buyOrSellOrBoth, dateStart, dateEnd, formClassArray, netTotalMin, netTotalMax, sharesMin, sharesMax, sharePriceMin, sharePriceMax, take, skip)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, forms)

}

func HandleDeepFilter(sortBy string, order string, buyOrSellOrBoth string, dateStart string, dateEnd string, formClass []string, netTotalMin float64, netTotalMax float64, sharesMin float64, sharesMax float64, sharePriceMin float64, sharePriceMax float64, take int64, skip int64) ([]structs.DB_DeltaForm, error) {

	buyOrSellFilter := bson.E{}
	if buyOrSellOrBoth == "Buy" {
		buyOrSellFilter = bson.E{Key: "buyOrSell", Value: "Buy"}
	} else if buyOrSellOrBoth == "Sell" {
		buyOrSellFilter = bson.E{Key: "buyOrSell", Value: "Sell"}
	}

	netTotalFilter := bson.E{}
	if netTotalMax < 100000000 {
		netTotalFilter = bson.E{Key: "netTotal", Value: bson.M{"$gte": netTotalMin, "$lte": netTotalMax}}
	}

	sharesFilter := bson.E{}
	if sharesMax < 100000000 {
		sharesFilter = bson.E{Key: "sharesTraded", Value: bson.M{"$gte": sharesMin, "$lte": sharesMax}}
	}

	sharePriceFilter := bson.E{}
	if sharePriceMax < 5000 {
		sharePriceFilter = bson.E{Key: "averagePricePerShare", Value: bson.M{"$gte": sharePriceMin, "$lte": sharePriceMax}}
	}
	// setup the filter
	filter := bson.D{
		{Key: "periodOfReport", Value: bson.D{{Key: "$gte", Value: dateStart}, {Key: "$lte", Value: dateEnd}}},
		{Key: "formClass", Value: bson.D{{Key: "$in", Value: formClass}}},
		sharePriceFilter,
		netTotalFilter,
		sharesFilter,
		buyOrSellFilter,
	}

	opts := options.Find().SetLimit(take).SetSkip(skip)

	queryOrder := -1
	if order == "asc" {
		queryOrder = 1
	}

	switch sortBy {
	case "periodOfReport":
		opts.SetSort(bson.D{{Key: "periodOfReport", Value: queryOrder}})
	case "netTotal":
		opts.SetSort(bson.D{{Key: "netTotal", Value: queryOrder}})
	case "sharesTraded":
		opts.SetSort(bson.D{{Key: "sharesTraded", Value: queryOrder}})
	case "averagePricePerShare":
		opts.SetSort(bson.D{{Key: "averagePricePerShare", Value: queryOrder}})
	}

	// run the query
	deltaForms, err := utils.DeltaFormFetch(filter, opts)

	if err != nil {
		return deltaForms, err
	}

	return deltaForms, nil
}
