package v1_handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"sync"

	"github.com/Matterhorn-Studios/insiderviz-forms-api/v1/models"
	"github.com/Matterhorn-Studios/insiderviz-forms-api/v1/v1_database"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/sync/errgroup"
)

func ClusterBuys(c *fiber.Ctx) error {
	// start with last month
	start_date := c.Query("start_date")
	end_date := c.Query("end_date")

	company_ciks, err := get_top_ten_clusters(start_date, end_date)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}

	// setup wait group
	var wg sync.WaitGroup
	wg.Add(len(company_ciks))

	// error channel
	errg := new(errgroup.Group)

	// get the company cards
	company_cards := make([]company_card, 10)
	for idx, company := range company_ciks {
		idx := idx
		company := company
		errg.Go(func() error {
			return func(idx int, company Cluster_Buy_Aggregation) error {
				wg.Done()
				company_cards[idx], err = get_company_card(company.Id, company.Ticker, company.Name, start_date, end_date)
				if err != nil {
					return err
				}
				return nil
			}(idx, company)
		})
	}

	if err := errg.Wait(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}

	// wait for all the company cards to be done
	wg.Wait()

	return c.JSON(company_cards)
}

type Cluster_Buy_Aggregation struct {
	Id        string               `bson:"_id"`
	Count     int                  `bson:"count"`
	Ticker    string               `bson:"ticker"`
	Name      string               `bson:"name"`
	Reporters []models.DB_Reporter `bson:"reporters"`
}

func get_top_ten_clusters(start_date string, end_date string) ([]Cluster_Buy_Aggregation, error) {
	ret_companies := make([]Cluster_Buy_Aggregation, 10)

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
			// add ticker
			{Key: "ticker", Value: bson.D{{Key: "$first", Value: "$issuer.issuerTicker"}}},
			// add name
			{Key: "name", Value: bson.D{{Key: "$first", Value: "$issuer.issuerName"}}},
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
		return ret_companies, err
	}

	if err = cursor.All(context.TODO(), &agg_result); err != nil {
		return ret_companies, err
	}

	// get the companies that are valid
	idx := 0
	comp_idx := 0
	for comp_idx < 10 && idx < len(agg_result) {
		cur_result := agg_result[idx]

		// ensure that there are at least three unique reporters
		reporters := make(map[string]bool)
		for _, reporter := range cur_result.Reporters {
			reporters[reporter.ReporterCik] = true
		}

		if len(reporters) >= 3 && valid_ticker(cur_result.Ticker) {
			ret_companies[comp_idx] = cur_result
			comp_idx++
		}

		idx++
	}

	return ret_companies, nil
}

func valid_ticker(ticker string) bool {
	if ticker == "" || ticker == "NONE" || ticker == "N/A" || ticker == "none" || ticker == "None" {
		return false
	}
	return true
}

type company_card struct {
	Company_Name        string               `json:"company_name"`
	Company_Cik         string               `json:"company_cik"`
	Company_Ticker      string               `json:"company_ticker"`
	Company_Stock_Data  []ext_stock_data     `json:"company_stock_data"`
	Company_Stock_Trade []company_card_trade `json:"company_stock_trade"`
}

type company_card_trade struct {
	Date      string  `json:"date"`
	Price     float32 `json:"price"`
	Shares    float32 `json:"shares"`
	Net_Total float32 `json:"net_total"`
	Name      string  `json:"name"`
}

func get_company_card(cik, ticker, name, start_date, end_date string) (company_card, error) {
	var ret company_card

	// add basic info
	ret.Company_Name = name
	ret.Company_Cik = cik
	ret.Company_Ticker = ticker

	// get the last month of stock data
	stock_data, err := get_stock_data(start_date, end_date, ticker)
	if err != nil {
		return ret, err
	}
	ret.Company_Stock_Data = stock_data

	// get the last month of trades
	trades, err := get_stock_trades(cik, start_date, end_date)
	if err != nil {
		return ret, err
	}
	ret.Company_Stock_Trade = trades

	return ret, nil
}

type ext_stock_data struct {
	Date           string  `json:"date"`
	Adjusted_Close float32 `json:"adjusted_close"`
}

func get_stock_data(start_date, end_date, ticker string) ([]ext_stock_data, error) {
	ret := make([]ext_stock_data, 0)

	// get stock data from ext api
	url := "https://eodhistoricaldata.com/api/eod/" + ticker + "?fmt=json&api_token=6288e433919037.08587703&order=d&from=" + start_date + "&to=" + end_date
	resp, err := http.Get(url)
	if err != nil {
		return ret, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ret, err
	}

	var ext_stock_data []ext_stock_data
	// parse the json
	err = json.Unmarshal(body, &ext_stock_data)
	if err != nil {
		return ret, err
	}

	return ext_stock_data, nil
}

func get_stock_trades(cik, start_date, end_date string) ([]company_card_trade, error) {
	ret := make([]company_card_trade, 0)

	// create a new DB session
	session, err := v1_database.GetNewSession()
	if err != nil {
		return ret, err
	}

	defer session.EndSession(context.Background())

	collection := session.Client().Database("insiderviz").Collection("DeltaForm")

	args := bson.D{
		{Key: "issuer.issuerCik", Value: cik},
		{Key: "periodOfReport", Value: bson.D{
			{Key: "$gte", Value: start_date},
		}},
		{Key: "periodOfReport", Value: bson.D{
			{Key: "$lte", Value: end_date},
		}},
		{Key: "buyOrSell", Value: "Buy"},
	}

	opts := options.Find().SetSort(bson.D{{Key: "periodOfReport", Value: -1}})

	// get the trades
	cursor, err := collection.Find(context.Background(), args, opts)
	if err != nil {
		return ret, err
	}

	var trades []models.DB_DeltaForm
	if err = cursor.All(context.Background(), &trades); err != nil {
		return ret, err
	}

	// pull out key info
	for _, trade := range trades {
		ret = append(ret, company_card_trade{
			Date:      trade.PeriodOfReport,
			Price:     trade.AveragePricePerShare,
			Shares:    trade.SharesTraded,
			Net_Total: trade.NetTotal,
			Name:      trade.Reporters[0].ReporterName,
		})
	}

	return ret, nil
}
