package single

import (
	"context"
	"math/rand"
	"net/http"
	"sync"

<<<<<<< HEAD
	"github.com/Matterhorn-Studios/insiderviz-forms-api/config"
=======
	iv_structs "github.com/Matterhorn-Studios/insiderviz-backend_structs"
	"github.com/Matterhorn-Studios/insiderviz-forms-api/database"
>>>>>>> 593ca62fdbab58a70df67f3c6af0b8ea92c171a6
	"github.com/Matterhorn-Studios/insiderviz-forms-api/v1/utils"
	"github.com/Matterhorn-Studios/insidervizforms/iv_models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Issuer(c *gin.Context) {
	cik := c.Param("cik")
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
		deltaForms, err = utils.DeltaFormFetch(filter, opts)
		if err != nil {
			deltaForms = []iv_models.DB_DeltaForm{}
		}
	}()

	var issuer iv_models.DB_Issuer_Doc
	go func() {
		defer wg.Done()
		var err error
		if includeGraph == "true" {
			issuer, err = utils.UpdateStockData(cik)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error1": err.Error()})
				return
			}
		} else {
			// get the issuer's information
			issuerCollection := database.GetCollection("Issuer")
			filter = bson.D{{Key: "cik", Value: cik}}
			opts := options.FindOne().SetProjection(bson.D{{Key: "stockData", Value: 0}})

			// do not include the stockData
			issuerInfo := issuerCollection.FindOne(context.TODO(), filter, opts)

			if issuerInfo.Err() != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error1": issuerInfo.Err().Error()})
				return
			}

			// unmarshal the issuer info
			if err = issuerInfo.Decode(&issuer); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error1": err.Error()})
				return
			}
		}
	}()
	wg.Wait()

	c.JSON(http.StatusOK, gin.H{"forms": deltaForms, "info": issuer})
}

func RandomIssuer(c *gin.Context) {
	// get the issuer collection
<<<<<<< HEAD
	issuerCollection := config.GetCollection(config.DB, "Issuer")
	var issuer iv_models.DB_Issuer_Doc
=======
	issuerCollection := database.GetCollection("Issuer")
	var issuer iv_structs.DB_Issuer_Doc
>>>>>>> 593ca62fdbab58a70df67f3c6af0b8ea92c171a6

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
					deltaForms, err := utils.DeltaFormFetch(filter, opts2)
					if err == nil {
						if len(deltaForms) > 0 && deltaForms[0].PeriodOfReport > "2021-01-01" {
							c.JSON(http.StatusOK, issuer)
							return
						}
					}
				}
			}
		}

	}
}

func IssuerGraph(c *gin.Context) {
	cik := c.Param("cik")

	// get the stock data
	stockData, err := utils.UpdateStockData(cik)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error1": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stockData)
}
