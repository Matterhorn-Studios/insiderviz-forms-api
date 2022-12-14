package v1_handlers

import (
	"context"
	"sync"

	"github.com/Matterhorn-Studios/insiderviz-forms-api/v1/v1_database"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Search(c *fiber.Ctx) error {
	// get the query
	query := c.Query("query")

	if query == "" {
		return c.JSON(fiber.Map{})
	}

	var wg sync.WaitGroup

	var issuers []IssuerRes
	var reporters []ReporterRes
	var err error
	wg.Add(2)
	go func() {
		defer wg.Done()
		issuers, err = searchIssuer(query, 10)
	}()
	go func() {
		defer wg.Done()
		reporters, err = searchReporter(query, 10)
	}()

	wg.Wait()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}

	sendItems := make([]SearchSend, 0)

	iIndex := 0
	rIndex := 0

	for len(sendItems) < 10 {
		curIssuer := issuers[iIndex]
		curReporter := reporters[rIndex]

		if curIssuer.Score < curReporter.Score {
			// take reporter
			rIndex++
			sendItems = append(sendItems, SearchSend{
				IsCongressman: curReporter.IsCongressman,
				IsIssuer:      false,
				Name:          curReporter.Name,
				Cik:           curReporter.Cik,
				Tickers:       make([]string, 0),
			})
		} else {
			// take issuer
			iIndex++
			sendItems = append(sendItems, SearchSend{
				IsCongressman: false,
				IsIssuer:      true,
				Name:          curIssuer.Name,
				Cik:           curIssuer.Cik,
				Tickers:       curIssuer.Tickers,
			})
		}
	}

	return c.JSON(sendItems)
}

type SearchSend struct {
	IsCongressman bool     `json:"isCongressman"`
	IsIssuer      bool     `json:"isIssuer"`
	Name          string   `json:"name"`
	Cik           string   `json:"cik"`
	Tickers       []string `json:"tickers"`
}

func searchIssuer(query string, limit int) ([]IssuerRes, error) {
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
			{
				Key: "score", Value: bson.D{
					{Key: "$meta", Value: "searchScore"},
				},
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

	limitFilter := bson.D{{Key: "$limit", Value: limit}}

	session, err := v1_database.GetNewSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(context.Background())

	// get the collection
	collection := session.Client().Database("insiderviz").Collection("Issuer")

	cursor, err := collection.Aggregate(context.TODO(), mongo.Pipeline{searchFilter, limitFilter, projectStage})
	var issuers []IssuerRes
	if err != nil {
		return issuers, err
	}

	if err = cursor.All(context.TODO(), &issuers); err != nil {
		return issuers, err
	}
	return issuers, nil
}

func searchReporter(query string, limit int) ([]ReporterRes, error) {
	projectStage := bson.D{
		{Key: "$project", Value: bson.D{
			{
				Key: "cik", Value: 1,
			},
			{
				Key: "name", Value: 1,
			},
			{
				Key: "isCongressman", Value: 1,
			},
			{
				Key: "score", Value: bson.D{
					{Key: "$meta", Value: "searchScore"},
				},
			},
		}},
	}

	// setup the filter
	searchFilter := bson.D{{Key: "$search", Value: bson.D{
		{Key: "index", Value: "reporter_name"},
		{Key: "text", Value: bson.D{
			{Key: "query", Value: query},
			{Key: "path", Value: bson.A{"name", "cik"}},
			{Key: "fuzzy", Value: bson.D{}},
		}},
	}}}

	limitFilter := bson.D{{Key: "$limit", Value: limit}}

	session, err := v1_database.GetNewSession()
	if err != nil {
		return nil, err
	}

	defer session.EndSession(context.Background())

	collection := session.Client().Database("insiderviz").Collection("Reporter")

	cursor, err := collection.Aggregate(context.TODO(), mongo.Pipeline{searchFilter, limitFilter, projectStage})
	var reporters []ReporterRes
	if err != nil {
		return reporters, err
	}

	if err = cursor.All(context.TODO(), &reporters); err != nil {
		return reporters, err
	}
	return reporters, nil
}

type ReporterRes struct {
	Cik           string  `json:"cik" bson:"cik"`
	Name          string  `json:"name" bson:"name"`
	IsCongressman bool    `json:"isCongressman" bson:"isCongressman"`
	Score         float64 `json:"score" bson:"score"`
}

type IssuerRes struct {
	Cik     string   `json:"cik" bson:"cik"`
	Name    string   `json:"name" bson:"name"`
	Tickers []string `json:"tickers" bson:"tickers"`
	Ein     string   `json:"ein" bson:"ein"`
	Score   float64  `json:"score" bson:"score"`
}
