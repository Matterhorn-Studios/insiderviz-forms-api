package utils

import (
	"context"

	"github.com/Matterhorn-Studios/insiderviz-forms-api/database"
	"github.com/Matterhorn-Studios/insidervizforms/iv_models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func CalcSectorHistory(startDate string) ([]iv_models.DB_Sector, error) {
	var results []iv_models.DB_Sector

	// setup the aggregate pipeline
	matchStage := bson.D{
		{Key: "$match", Value: bson.D{
			{Key: "periodOfReport", Value: bson.D{{Key: "$gte", Value: startDate}}},
			{Key: "issuer.issuerSector", Value: bson.D{{Key: "$in", Value: bson.A{"Financial", "Real Estate", "Healthcare", "Consumer Defensive", "Fund", "Energy", "Basic Materials", "Industrials", "Technology", "Utilities", "Consumer Cyclical", "Communication Services", "Financial Services"}}}},
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
				Key: "netTotal", Value: 1,
			},
			{
				Key: "periodOfReport", Value: 1,
			},
			{
				Key: "issuer", Value: 1,
			},
		}},
	}

	// group by sector
	groupStage := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{
				{Key: "date", Value: "$periodOfReport"},
				{Key: "name", Value: "$issuer.issuerSector"},
			}},
			{Key: "totalBought", Value: bson.D{
				{Key: "$sum", Value: "$SumBuys"},
			}},
			{Key: "totalSold", Value: bson.D{
				{Key: "$sum", Value: "$SumSells"},
			}},
			{Key: "total", Value: bson.D{
				{Key: "$sum", Value: "$netTotal"},
			}},
		}},
	}

	// sort
	orderStage := bson.D{
		{Key: "$sort", Value: bson.D{
			{Key: "_id.date", Value: -1},
		}},
	}

	groupStage3 := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{
				{Key: "name", Value: "$_id.name"},
			}},
			{Key: "historicalData", Value: bson.D{
				{Key: "$push", Value: bson.D{
					{Key: "date", Value: "$_id.date"},
					{Key: "totalBought", Value: "$totalBought"},
					{Key: "totalSold", Value: "$totalSold"},
					{Key: "total", Value: "$total"},
					{Key: "buyOrSell", Value: bson.D{
						{Key: "$cond", Value: bson.A{
							bson.D{{
								Key: "$gte", Value: bson.A{"$totalBought", "$totalSold"},
							}},
							"Buy", "Sell",
						}},
					}},
				}},
			}},
		}},
	}

	cursor, err := database.GetCollection("DeltaForm").Aggregate(context.TODO(),
		mongo.Pipeline{
			matchStage,
			projectStage,
			groupStage,
			orderStage,
			groupStage3,
		})

	if err != nil {
		return results, err
	}

	if err = cursor.All(context.TODO(), &results); err != nil {
		return results, err
	}

	return results, nil
}
