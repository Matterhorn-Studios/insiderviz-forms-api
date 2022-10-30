package v1_handlers

import (
	"context"
	"fmt"

	"github.com/Matterhorn-Studios/insiderviz-forms-api/v1/models"
	"github.com/Matterhorn-Studios/insiderviz-forms-api/v1/v1_database"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func ClusterBuys(c *fiber.Ctx) error {
	// start with last month
	start_date := "2022-10-01"
	end_date := "2022-10-31"

	match_stage := bson.D{
		{Key: "$match", Value: bson.D{
			{Key: "periodOfReport", Value: bson.D{{Key: "$gte", Value: start_date}}},
			{Key: "periodOfReport", Value: bson.D{{Key: "$lte", Value: end_date}}},
			{Key: "buyOrSell", Value: "Buy"},
		}},
	}

	project_stage := bson.D{
		{Key: "$project", Value: bson.D{
			{
				Key: "issuer", Value: 1,
			},
			{
				Key: "reporters", Value: 1,
			},
		}},
	}

	group_stage := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$issuer.issuerCik"},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
			// add first element of reporters array
			{Key: "reporters", Value: bson.D{{Key: "$push", Value: bson.D{{Key: "$arrayElemAt", Value: bson.A{"$reporters", 0}}}}}},
		}},
	}

	order_stage := bson.D{
		{Key: "$sort", Value: bson.D{{
			Key: "count", Value: -1,
		}}},
	}

	limit_stage := bson.D{{Key: "$limit", Value: 25}}

	var agg_result []Cluster_Buy_Aggregation

	cursor, err := v1_database.GetCollection("DeltaForm").Aggregate(context.TODO(), mongo.Pipeline{match_stage, project_stage, group_stage, order_stage, limit_stage})
	if err != nil {
		return nil
	}

	if err = cursor.All(context.TODO(), &agg_result); err != nil {
		return nil
	}

	// get the companies that are valid
	var companies []string
	idx := 0
	for len(companies) < 10 && idx < len(agg_result) {
		fmt.Println("checking company", agg_result[idx].Id)
		cur_result := agg_result[idx]

		// ensure that there are at least three unique reporters
		reporters := make(map[string]bool)
		for _, reporter := range cur_result.Reporters {
			reporters[reporter.ReporterCik] = true
		}

		if len(reporters) >= 3 {
			companies = append(companies, cur_result.Id)
		}

		idx++
	}

	return c.JSON(companies)
}

type Cluster_Buy_Aggregation struct {
	Id        string               `bson:"_id"`
	Count     int                  `bson:"count"`
	Reporters []models.DB_Reporter `bson:"reporters"`
}
