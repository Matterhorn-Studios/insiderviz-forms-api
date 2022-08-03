package utils

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/Matterhorn-Studios/insiderviz-forms-api/config"
	"github.com/Matterhorn-Studios/insiderviz-forms-api/v1/structs"
	"go.mongodb.org/mongo-driver/bson"
)

func UpdateStockData(cik string) (structs.DB_Issuer_Doc, error) {
	var issuer structs.DB_Issuer_Doc

	// get the current issuer info
	issuerCollection := config.GetCollection(config.DB, "Issuer")
	filter := bson.D{{Key: "cik", Value: cik}}

	issuerInfo := issuerCollection.FindOne(context.TODO(), filter)

	if issuerInfo.Err() != nil {
		return issuer, issuerInfo.Err()
	}

	// unmarshal the issuer info
	if err := issuerInfo.Decode(&issuer); err != nil {
		return issuer, err
	}

	// get today's date in the format YYYY-MM-DD
	today := time.Now().Format("2006-01-02")

	// check the length of the stock data
	if len(issuer.StockData) == 0 && len(issuer.Tickers) > 0 {
		// fetch new data
		url := "https://eodhistoricaldata.com/api/eod/" + issuer.Tickers[0] + "?fmt=json&api_token=6288e433919037.08587703&order=d&from=2016-01-01"

		// create the http request
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return issuer, err
		}

		// create the http client
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil || resp.StatusCode != 200 {
			return issuer, err
		}

		// convert response body to byte array
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return issuer, err
		}

		// parse
		var apiRes []stockDataAPIResponse
		err = json.Unmarshal(bodyBytes, &apiRes)
		if err != nil {
			return issuer, err
		}

		// save it to the issuer
		for _, day := range apiRes {
			issuer.StockData = append(issuer.StockData, structs.StockData{
				Date:   day.Date,
				Close:  day.Close,
				Volume: day.Volume,
			})
		}

		// update the issuer
		filter := bson.D{{Key: "cik", Value: cik}}
		_, err = issuerCollection.UpdateOne(context.TODO(), filter, bson.D{{Key: "$set", Value: issuer}})
		if err != nil {
			return issuer, err
		}
	} else if len(issuer.Tickers) > 0 && issuer.StockData[0].Date != today {
		// need to update the stock data
		startDate := issuer.StockData[0].Date
		url := "https://eodhistoricaldata.com/api/eod/" + issuer.Tickers[0] + "?fmt=json&api_token=6288e433919037.08587703&order=d&from=" + startDate

		// create the http request
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return issuer, err
		}

		// create the http client
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil || resp.StatusCode != 200 {
			return issuer, err
		}

		// convert response body to byte array
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return issuer, err
		}

		// parse
		var apiRes []stockDataAPIResponse
		err = json.Unmarshal(bodyBytes, &apiRes)
		if err != nil {
			return issuer, err
		}

		for _, day := range apiRes {
			if day.Date != startDate {
				temp := []structs.StockData{{
					Date:   day.Date,
					Close:  day.Close,
					Volume: day.Volume,
				}}
				issuer.StockData = append(temp, issuer.StockData...)
			}
		}

		// update the issuer
		filter := bson.D{{Key: "cik", Value: cik}}
		_, err = issuerCollection.UpdateOne(context.TODO(), filter, bson.D{{Key: "$set", Value: issuer}})
		if err != nil {
			return issuer, err
		}

	}

	return issuer, nil

}

type stockDataAPIResponse struct {
	Date           string  `json:"date"`
	Open           float64 `json:"open"`
	High           float64 `json:"high"`
	Low            float64 `json:"low"`
	Close          float64 `json:"close"`
	Adjusted_Close float64 `json:"adjusted_close"`
	Volume         int     `json:"volume"`
}
