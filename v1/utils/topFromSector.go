package utils

import (
	"context"

	"github.com/Matterhorn-Studios/insiderviz-forms-api/database"
	"github.com/Matterhorn-Studios/insidervizforms/iv_models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func TopFromSector(startDate string, sector string) (iv_models.Top_From_Sector, error) {
	var result iv_models.Top_From_Sector
	var aggResult []iv_models.Top_From_Sector_Entry

	matchStage := bson.D{
		{Key: "$match", Value: bson.D{
			{Key: "periodOfReport", Value: bson.D{{Key: "$gte", Value: startDate}}},
			{Key: "issuer.issuerSector", Value: bson.D{{Key: "$eq", Value: sector}}},
		}},
	}

	projectStage := bson.D{
		{Key: "$project", Value: bson.D{
			{
				Key: "netTotal", Value: 1,
			},
			{
				Key: "issuer", Value: 1,
			},
		}},
	}

	groupStage := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$issuer.issuerCik"},
			{Key: "name", Value: bson.D{
				{Key: "$first", Value: "$issuer.issuerName"},
			}},
			{Key: "industry", Value: bson.D{{Key: "$first", Value: "$issuer.issuerSector"}}},
			{Key: "ticker", Value: bson.D{{Key: "$first", Value: "$issuer.issuerTicker"}}},
			{Key: "tradeVolume", Value: bson.D{{Key: "$sum", Value: "$netTotal"}}},
		}},
	}

	orderStage := bson.D{
		{Key: "$sort", Value: bson.D{{
			Key: "tradeVolume", Value: -1,
		}}},
	}

	limitStage := bson.D{{Key: "$limit", Value: 5}}

	cursor, err := database.GetCollection("DeltaForm").Aggregate(context.TODO(), mongo.Pipeline{matchStage, projectStage, groupStage, orderStage, limitStage})
	if err != nil {
		return result, err
	}

	if err = cursor.All(context.TODO(), &aggResult); err != nil {
		return result, err
	}

	result.Sector = sector
	result.Companies = aggResult

	return result, nil
}
