package aggregation

import (
	"net/http"
	"sync"

	top "github.com/Matterhorn-Studios/insiderviz-forms-api/v1/utils/TopCompanies"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Top(c *gin.Context) {
	// get the start date from the query
	startDate := c.Query("startDate")

	var err error
	var mostBoughtInsider []primitive.M
	var mostSoldInsider []primitive.M
	var mostBoughtCongress []primitive.M
	var mostSoldCongress []primitive.M
	var wg sync.WaitGroup
	wg.Add(4)

	go func() {
		defer wg.Done()
		mostBoughtInsider, err = top.MostBought(startDate, "Insider")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	}()

	go func() {
		defer wg.Done()
		mostSoldInsider, err = top.MostSold(startDate, "Insider")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	}()

	go func() {
		defer wg.Done()
		mostBoughtCongress, err = top.MostBought(startDate, "Congress")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	}()

	go func() {
		defer wg.Done()
		mostSoldCongress, err = top.MostSold(startDate, "Congress")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	}()

	wg.Wait()

	returnObj := bson.M{
		"mostBoughtInsider":  mostBoughtInsider,
		"mostSoldInsider":    mostSoldInsider,
		"mostBoughtCongress": mostBoughtCongress,
		"mostSoldCongress":   mostSoldCongress,
	}

	c.JSON(http.StatusOK, returnObj)
}
