package top

import (
	"context"

	"github.com/Matterhorn-Studios/insiderviz-forms-api/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func MostSold(startDate string, formClass string) ([]bson.M, error) {
	var results []bson.M

	// setup the aggregate pipeline
	matchStage := bson.D{
		{Key: "$match", Value: bson.D{
			{Key: "periodOfReport", Value: bson.D{{Key: "$gte", Value: startDate}}},
			{Key: "formClass", Value: bson.D{{Key: "$eq", Value: formClass}}},
		}},
	}

	// sum the buys
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
			{Key: "amount", Value: bson.D{
				{Key: "$sum", Value: "$SumSells"},
			}},
			{Key: "name", Value: bson.D{
				{Key: "$first", Value: "$issuer.issuerName"},
			}},
		}},
	}

	// order
	orderState := bson.D{
		{Key: "$sort", Value: bson.D{{Key: "amount", Value: -1}}},
	}

	// limit
	limitStage := bson.D{{Key: "$limit", Value: 10}}

	// run the aggregate on delta form
	cursor, err := database.GetCollection("DeltaForm").Aggregate(context.TODO(), mongo.Pipeline{matchStage, projectStage, groupStage, orderState, limitStage})
	if err != nil {
		return results, err
	}

	// display the results
	if err = cursor.All(context.TODO(), &results); err != nil {
		return results, err
	}

	return results, nil
}
