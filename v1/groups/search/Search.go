package search

import (
	"context"
	"net/http"

	"github.com/Matterhorn-Studios/insiderviz-forms-api/config"
	"github.com/Matterhorn-Studios/insiderviz-forms-api/v1/structs"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Search(c *gin.Context) {
	// get the query
	query := c.Query("query")

	if query == "" {
		c.JSON(http.StatusOK, make([]string, 0))
		return
	}

	projectStage := bson.D{
		{Key: "$project", Value: bson.D{
			{
				Key: "cik", Value: 1,
			},
			{
				Key: "name", Value: 1,
			},
			{
				Key: "tickers", Value: 1,
			},
			{
				Key: "ein", Value: 1,
			},
		}},
	}

	// setup the filter
	searchFilter := bson.D{{Key: "$search", Value: bson.D{
		{Key: "index", Value: "issuer_name"},
		{Key: "text", Value: bson.D{
			{Key: "query", Value: query},
			{Key: "path", Value: bson.A{"name", "cik", "ein", "tickers"}},
			{Key: "fuzzy", Value: bson.D{}},
		}},
	}}}

	limitFilter := bson.D{{Key: "$limit", Value: 10}}

	cursor, err := config.GetCollection(config.DB, "Issuer").Aggregate(context.TODO(), mongo.Pipeline{searchFilter, projectStage, limitFilter})
	if err != nil {
		panic(err)
	}

	var issuers []structs.DB_Issuer_Doc
	if err = cursor.All(context.TODO(), &issuers); err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, issuers)
}
