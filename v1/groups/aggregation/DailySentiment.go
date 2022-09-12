package aggregation

import (
	"context"
	"net/http"

	"github.com/Matterhorn-Studios/insiderviz-forms-api/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func DailySentiment(c *gin.Context) {
	// get the start date from query
	startDate := c.Query("startDate")

	// setup the aggregate pipeline
	matchStage := bson.D{
		{Key: "$match", Value: bson.D{
			{Key: "periodOfReport", Value: bson.D{{Key: "$gte", Value: startDate}}},
		}},
	}

	projectStage := bson.D{
		{Key: "$project", Value: bson.D{
			{
				Key: "SumSells", Value: bson.D{
					{Key: "$cond", Value: bson.A{
						bson.D{{
							Key: "$eq", Value: bson.A{"$buyOrSell", "Sell"},
						}},
						"$netTotal", 0,
					}},
				},
			},
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
				Key: "CountSells", Value: bson.D{
					{Key: "$cond", Value: bson.A{
						bson.D{{
							Key: "$eq", Value: bson.A{"$buyOrSell", "Sell"},
						}},
						1, 0,
					}},
				},
			},
			{
				Key: "CountBuys", Value: bson.D{
					{Key: "$cond", Value: bson.A{
						bson.D{{
							Key: "$eq", Value: bson.A{"$buyOrSell", "Buy"},
						}},
						1, 0,
					}},
				},
			},
			{
				Key: "periodOfReport", Value: 1,
			},
		}},
	}

	groupStage := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$periodOfReport"},
			{Key: "sellAmount", Value: bson.D{
				{Key: "$sum", Value: "$SumSells"},
			}},
			{Key: "buyAmount", Value: bson.D{
				{Key: "$sum", Value: "$SumBuys"},
			}},
			{Key: "sellCount", Value: bson.D{
				{Key: "$sum", Value: "$CountSells"},
			}},
			{Key: "buyCount", Value: bson.D{
				{Key: "$sum", Value: "$CountBuys"},
			}},
		}},
	}

	orderState := bson.D{
		{Key: "$sort", Value: bson.D{{Key: "_id", Value: -1}}},
	}

	// run the aggregate
	cursor, err := database.GetCollection("DeltaForm").Aggregate(context.TODO(), mongo.Pipeline{matchStage, projectStage, groupStage, orderState})
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
