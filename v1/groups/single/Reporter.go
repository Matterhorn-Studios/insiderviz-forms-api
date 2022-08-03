package single

import (
	"context"
	"math/rand"
	"net/http"

	"github.com/Matterhorn-Studios/insiderviz-forms-api/config"
	"github.com/Matterhorn-Studios/insiderviz-forms-api/v1/structs"
	"github.com/Matterhorn-Studios/insiderviz-forms-api/v1/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Reporter(c *gin.Context) {
	cik := c.Param("cik")

	filter := bson.D{{Key: "reporters.reporterCik", Value: cik}}
	opts := options.Find().SetSort(bson.D{{Key: "periodOfReport", Value: -1}})

	deltaForms, err := utils.DeltaFormFetch(filter, opts)

	if err != nil {
		deltaForms = []structs.DB_DeltaForm{}
	}

	// get the reporter's information
	reporterCollection := config.GetCollection(config.DB, "Reporter")
	filter = bson.D{{Key: "cik", Value: cik}}

	issuerInfo := reporterCollection.FindOne(context.TODO(), filter)

	if issuerInfo.Err() != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error1": issuerInfo.Err().Error()})
		return
	}

	// unmarshal the reporter info
	var reporter bson.M
	if err = issuerInfo.Decode(&reporter); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error1": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"forms": deltaForms, "info": reporter})
}

func RandomReporter(c *gin.Context) {
	// get the reporter collection
	reporterCollection := config.GetCollection(config.DB, "Reporter")
	var reporter structs.DB_Reporter_Doc

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
				deltaForms, err := utils.DeltaFormFetch(filter, opts2)
				if err == nil {
					if len(deltaForms) > 0 && deltaForms[0].PeriodOfReport > "2021-01-01" {
						c.JSON(http.StatusOK, reporter)
						return
					}
				}
			}
		}

	}
}
