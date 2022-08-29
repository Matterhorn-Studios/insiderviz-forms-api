package single

import (
	"context"
	"math/rand"
	"net/http"
	"sort"
	"strconv"

	"github.com/Matterhorn-Studios/insiderviz-forms-api/config"
	"github.com/Matterhorn-Studios/insiderviz-forms-api/v1/structs"
	"github.com/Matterhorn-Studios/insiderviz-forms-api/v1/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func LatestThirteenF(c *gin.Context) {
	cik := c.Param("cik")

	rawOffset := c.Query("offset")
	offset := 0
	if rawOffset != "" {
		var err error
		offset, err = strconv.Atoi(rawOffset)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	thirteenFCollection := config.GetCollection(config.DB, "13F")

	filter := bson.D{{Key: "cik", Value: cik}}
	opts := options.FindOne().SetSort(bson.D{{Key: "periodOfReport", Value: -1}}).SetSkip(int64(offset))

	form := thirteenFCollection.FindOne(context.TODO(), filter, opts)

	if form.Err() != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": form.Err().Error()})
		return
	} else {
		var thirteenF structs.DB_Form13F_Base
		err := form.Decode(&thirteenF)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
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
				thirteenF.Holdings = append(thirteenF.Holdings, structs.DB_Form13F_Holding{
					Name:     "Other",
					NetTotal: float32(otherTotal),
					Shares:   float32(otherShares),
					Cik:      "",
				})
			}

			c.JSON(http.StatusOK, thirteenF)
		}
	}

}

func Reporter(c *gin.Context) {
	cik := c.Param("cik")

	filter := bson.D{{Key: "reporters.reporterCik", Value: cik}, {
		Key: "$or", Value: bson.A{
			bson.D{{Key: "formClass", Value: "Insider"}},
			bson.D{{Key: "formClass", Value: "Congress"}},
		},
	}}
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
