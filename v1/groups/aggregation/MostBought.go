package aggregation

import (
	"context"
	"net/http"

	"github.com/Matterhorn-Studios/insiderviz-forms-api/config"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func MostBought(c *gin.Context) {
	// get the start date from the query
	startDate := c.Query("startDate")

	// setup the aggregate pipeline
	matchStage := bson.D{
		{Key: "$match", Value: bson.D{
			{Key: "periodOfReport", Value: bson.D{{Key: "$gte", Value: startDate}}},
		}},
	}

	// sum the buys
	projectStage := bson.D{
		{Key: "$project", Value: bson.D{
			{
				Key: "SumBuys", Value: bson.D{
					{Key: "$cond", Value: bson.A{
						bson.D{{
							Key: "$eq", Value: bson.A{"$buyOrSell", "Buy"},
						}},
						"$netTotal", 0,
					}},
				},
			},
			{
				Key: "periodOfReport", Value: 1,
			},
			{
				Key: "issuer", Value: 1,
			},
		}},
	}

	// group by company
	groupStage := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$issuer.issuerCik"},
			{Key: "buyAmount", Value: bson.D{
				{Key: "$sum", Value: "$SumBuys"},
			}},
			{Key: "name", Value: bson.D{
				{Key: "$first", Value: "$issuer.issuerName"},
			}},
		}},
	}

	// order
	orderState := bson.D{
		{Key: "$sort", Value: bson.D{{Key: "buyAmount", Value: -1}}},
	}

	// limit
	limitStage := bson.D{{Key: "$limit", Value: 10}}

	// run the aggregate on delta form
	cursor, err := config.GetCollection(config.DB, "DeltaForm").Aggregate(context.TODO(), mongo.Pipeline{matchStage, projectStage, groupStage, orderState, limitStage})
	if err != nil {
		panic(err)
	}

	// display the results
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, results)
}
