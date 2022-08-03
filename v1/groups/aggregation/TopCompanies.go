package aggregation

import (
	"net/http"

	top "github.com/Matterhorn-Studios/insiderviz-forms-api/v1/utils/TopCompanies"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func Top(c *gin.Context) {
	// get the start date from the query
	startDate := c.Query("startDate")

	mostBoughtInsider, err := top.MostBought(startDate, "Insider")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	mostSoldInsider, err := top.MostSold(startDate, "Insider")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	mostBoughtCongress, err := top.MostBought(startDate, "Congress")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	mostSoldCongress, err := top.MostSold(startDate, "Congress")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	returnObj := bson.M{
		"mostBoughtInsider":  mostBoughtInsider,
		"mostSoldInsider":    mostSoldInsider,
		"mostBoughtCongress": mostBoughtCongress,
		"mostSoldCongress":   mostSoldCongress,
	}

	c.JSON(http.StatusOK, returnObj)
}
