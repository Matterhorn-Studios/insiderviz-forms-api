package lib

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Matterhorn-Studios/insiderviz-forms-api/v1/v1_database"
	"github.com/Matterhorn-Studios/insidervizforms/iv_models"
	"go.mongodb.org/mongo-driver/bson"
)

func UpdateStockData(cik string) (iv_models.DB_Issuer_Doc, error) {
	var issuer iv_models.DB_Issuer_Doc

	// get the current issuer info
	issuerCollection := v1_database.GetCollection("Issuer")
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
	if (len(issuer.StockData) == 0 && len(issuer.Tickers) > 0) || (!issuer.StockDataSplit && len(issuer.Tickers) > 0) {

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
		stockDataBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return issuer, err
		}

		// get the stock split data
		url = "https://eodhistoricaldata.com/api/splits/" + issuer.Tickers[0] + "?fmt=json&api_token=6288e433919037.08587703&order=d&from=2016-01-01"

		// create the http request
		req, err = http.NewRequest("GET", url, nil)
		if err != nil {
			return issuer, err
		}

		// create the http client
		resp, err = client.Do(req)
		if err != nil || resp.StatusCode != 200 {
			return issuer, err
		}

		// convert response body to byte array
		stockSplitBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return issuer, err
		}

		// parse
		var apiRes []stockDataAPIResponse
		err = json.Unmarshal(stockDataBytes, &apiRes)
		if err != nil {
			return issuer, err
		}

		var apiSplitRes []stockSplitAPIResponse
		err = json.Unmarshal(stockSplitBytes, &apiSplitRes)
		if err != nil {
			return issuer, err
		}

		// save current index for split
		curSplitIndex := len(apiSplitRes) - 1
		curSplit := 1.0
		if curSplitIndex >= 0 {
			splitStr := strings.Split(apiSplitRes[curSplitIndex].Split, "/")
			top, _ := strconv.ParseFloat(splitStr[0], 64)
			bot, _ := strconv.ParseFloat(splitStr[1], 64)
			curSplit *= top / bot
		}

		issuer.StockData = make([]iv_models.StockData, 0)

		for _, day := range apiRes {
			// check the split
			// ensure there is a split left
			if curSplitIndex >= 0 {
				if apiSplitRes[curSplitIndex].Date > day.Date {
					// make sure it is the most recent
					if curSplitIndex-1 >= 0 {
						if apiSplitRes[curSplitIndex-1].Date > day.Date {
							curSplitIndex--
							splitStr := strings.Split(apiSplitRes[curSplitIndex].Split, "/")
							top, _ := strconv.ParseFloat(splitStr[0], 64)
							bot, _ := strconv.ParseFloat(splitStr[1], 64)
							curSplit *= top / bot
						}
					}

					// apply the split
					day.Close /= curSplit
				}
				temp := iv_models.StockData{
					Date:   day.Date,
					Close:  day.Close,
					Volume: day.Volume,
				}
				issuer.StockData = append(issuer.StockData, temp)
			} else {
				temp := iv_models.StockData{
					Date:   day.Date,
					Close:  day.Close,
					Volume: day.Volume,
				}
				issuer.StockData = append(issuer.StockData, temp)
			}
		}

		// update the issuer
		issuer.StockDataSplit = true
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
		stockDataBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return issuer, err
		}

		// get the stock split data
		url = "https://eodhistoricaldata.com/api/splits/" + issuer.Tickers[0] + "?fmt=json&api_token=6288e433919037.08587703&order=d&from=" + startDate

		// create the http request
		req, err = http.NewRequest("GET", url, nil)
		if err != nil {
			return issuer, err
		}

		// create the http client
		resp, err = client.Do(req)
		if err != nil || resp.StatusCode != 200 {
			return issuer, err
		}

		// convert response body to byte array
		stockSplitBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return issuer, err
		}

		// parse
		var apiRes []stockDataAPIResponse
		err = json.Unmarshal(stockDataBytes, &apiRes)
		if err != nil {
			return issuer, err
		}

		var apiSplitRes []stockSplitAPIResponse
		err = json.Unmarshal(stockSplitBytes, &apiSplitRes)
		if err != nil {
			return issuer, err
		}

		// save current index for split
		curSplitIndex := len(apiSplitRes) - 1
		curSplit := 1.0
		if curSplitIndex >= 0 {
			splitStr := strings.Split(apiSplitRes[curSplitIndex].Split, "/")
			top, _ := strconv.ParseFloat(splitStr[0], 64)
			bot, _ := strconv.ParseFloat(splitStr[1], 64)

			curSplit *= top / bot
		}

		addGroup := make([]iv_models.StockData, 0)

		for _, day := range apiRes {
			if day.Date != startDate {
				// check the split
				// ensure there is a split left
				if curSplitIndex >= 0 {
					if apiSplitRes[curSplitIndex].Date > day.Date {
						// make sure it is the most recent
						if curSplitIndex-1 >= 0 {
							if apiSplitRes[curSplitIndex-1].Date > day.Date {
								curSplitIndex--
								splitStr := strings.Split(apiSplitRes[curSplitIndex].Split, "/")
								top, _ := strconv.ParseFloat(splitStr[0], 64)
								bot, _ := strconv.ParseFloat(splitStr[1], 64)
								curSplit *= top / bot
							}
						}

						// apply the split
						day.Close /= curSplit
					}
				}
				temp := iv_models.StockData{
					Date:   day.Date,
					Close:  day.Close,
					Volume: day.Volume,
				}
				addGroup = append(addGroup, temp)
			}
		}

		// prepend to the issuer
		issuer.StockData = append(addGroup, issuer.StockData...)

		// update the issuer
		filter := bson.D{{Key: "cik", Value: cik}}
		_, err = issuerCollection.UpdateOne(context.TODO(), filter, bson.D{{Key: "$set", Value: issuer}})
		if err != nil {
			return issuer, err
		}

	}

	return issuer, nil

}

type stockSplitAPIResponse struct {
	Date  string `json:"date"`
	Split string `json:"split"`
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
