package v1_handlers

import (
	"context"
	"sync"

	"github.com/Matterhorn-Studios/insiderviz-forms-api/v1/v1_database"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SendConfigReporter struct {
	Name                string `json:"name"`
	CIK                 string `json:"cik"`
	IsCongressman       bool   `json:"isCongressman"`
	MostRecentTradeDate string `json:"mostRecentTradeDate"`
	NumberOfTrades      int    `json:"numberOfTrades"`
}

func SearchConfigReporter(c *fiber.Ctx) error {
	query := c.Query("query")

	if query == "" {
		return c.JSON(fiber.Map{})
	}

	reporters, err := searchReporter(query, 5)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}

	sendReporters := make([]SendConfigReporter, len(reporters))

	var wg sync.WaitGroup

	for i := 0; i < len(sendReporters); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			reporter := reporters[i]
			lastDate := ""
			numberOfTrades := 0

			var innerWg sync.WaitGroup
			innerWg.Add(2)

			// get most recent trade date
			go func(date *string) {
				defer innerWg.Done()

				session, err := v1_database.GetNewSession()
				if err == nil {
					defer session.EndSession(context.Background())

					deltaCollection := session.Client().Database("insiderviz").Collection("DeltaForm")

					filter := bson.M{"reporters.reporterCik": reporter.Cik}
					opts := options.FindOne().SetProjection(bson.D{{Key: "periodOfReport", Value: 1}})

					res := deltaCollection.FindOne(context.Background(), filter, opts)
					if res.Err() == nil {
						var deltaRes struct {
							PeriodOfReport string `bson:"periodOfReport"`
						}

						if err := res.Decode(&deltaRes); err == nil {
							*date = deltaRes.PeriodOfReport
						}
					}
				}
			}(&lastDate)

			// get number of trades
			go func(number *int) {
				defer innerWg.Done()

				session, err := v1_database.GetNewSession()
				if err == nil {
					defer session.EndSession(context.Background())

					deltaCollection := session.Client().Database("insiderviz").Collection("DeltaForm")

					filter := bson.M{"reporters.reporterCik": reporter.Cik}
					count, err := deltaCollection.CountDocuments(context.Background(), filter)
					if err == nil {
						*number = int(count)
					}
				}
			}(&numberOfTrades)

			innerWg.Wait()

			sendReporters[i] = SendConfigReporter{
				Name:                reporter.Name,
				CIK:                 reporter.Cik,
				IsCongressman:       reporter.IsCongressman,
				MostRecentTradeDate: lastDate,
				NumberOfTrades:      numberOfTrades,
			}
		}(i)
	}

	wg.Wait()

	return c.JSON(sendReporters)

}

type SendConfigIssuer struct {
	Name                 string  `json:"name"`
	Ticker               string  `json:"ticker"`
	MostRecentTradeDate  string  `json:"mostRecentTradeDate"`
	MostRecentStockPrice float64 `json:"mostRecentStockPrice"`
	CIK                  string  `json:"cik"`
}

func SearchConfigIssuer(c *fiber.Ctx) error {
	// get the query
	query := c.Query("query")

	if query == "" {
		return c.JSON(fiber.Map{})
	}

	issuers, err := searchIssuer(query, 5)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}

	// get most recent trade date for each issuer
	sendIssuers := make([]SendConfigIssuer, len(issuers))

	var wg sync.WaitGroup

	for i := 0; i < len(sendIssuers); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			issuer := issuers[i]
			stockPrice := 0.0
			lastDate := ""

			var innerWg sync.WaitGroup

			innerWg.Add(2)

			// get the most recent stock price
			go func(price *float64) {
				defer innerWg.Done()

				session, err := v1_database.GetNewSession()
				if err == nil {
					defer session.EndSession(context.Background())

					issuerCollection := session.Client().Database("insiderviz").Collection("Issuer")

					filter := bson.M{"cik": issuer.Cik}
					opts := options.FindOne().SetProjection(bson.D{{Key: "stockData", Value: 1}})

					res := issuerCollection.FindOne(context.Background(), filter, opts)

					if res.Err() == nil {
						var issuerRes struct {
							StockData []struct {
								Price float64 `bson:"close"`
							} `bson:"stockData"`
						}

						if err := res.Decode(&issuerRes); err == nil {
							if len(issuerRes.StockData) > 0 {
								*price = issuerRes.StockData[0].Price
							}
						}
					}
				}
			}(&stockPrice)

			// get most recent trade date
			go func(date *string) {
				defer innerWg.Done()

				session, err := v1_database.GetNewSession()
				if err == nil {
					defer session.EndSession(context.Background())

					deltaCollection := session.Client().Database("insiderviz").Collection("DeltaForm")

					filter := bson.M{"issuer.issuerCik": issuer.Cik}
					opts := options.FindOne().SetProjection(bson.D{{Key: "periodOfReport", Value: 1}})

					res := deltaCollection.FindOne(context.Background(), filter, opts)
					if res.Err() == nil {
						var deltaRes struct {
							PeriodOfReport string `bson:"periodOfReport"`
						}

						if err := res.Decode(&deltaRes); err == nil {
							*date = deltaRes.PeriodOfReport
						}
					}
				}

			}(&lastDate)

			innerWg.Wait()

			ticker := ""
			if len(issuer.Tickers) > 0 {
				ticker = issuer.Tickers[0]
			}

			sendIssuers[i] = SendConfigIssuer{
				Name:                 issuer.Name,
				Ticker:               ticker,
				MostRecentTradeDate:  lastDate,
				MostRecentStockPrice: stockPrice,
				CIK:                  issuer.Cik,
			}
		}(i)
	}

	wg.Wait()

	return c.JSON(sendIssuers)
}
