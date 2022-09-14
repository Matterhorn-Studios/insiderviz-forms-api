package v1_handlers

import (
	"github.com/Matterhorn-Studios/insiderviz-forms-api/v1/lib"
	"github.com/Matterhorn-Studios/insidervizforms/iv_models"
	"github.com/gofiber/fiber/v2"
)

func SectorHistory(c *fiber.Ctx) error {
	// get the start datefrom the query
	startDate := c.Query("startDate")

	data, err := lib.CalcSectorHistory(startDate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}

	sendData := make([]SectorHistoryData, 0)

	for _, v := range data {
		addSectorToData(&sendData, v)
	}

	return c.JSON(sendData)
}

func addSectorToData(data *[]SectorHistoryData, sector iv_models.DB_Sector) {
	// ASSUME data IS IN ORDER
	for _, histItem := range sector.HistoricalData {
		found := false
		for i := range *data {
			curDataPt := (*data)[i]

			if curDataPt.Date == histItem.Date {
				// append to already existing entry
				appendCorrectSector(&(*data)[i], sector.Id.Name, histItem.Total)
				found = true
				break
			} else if curDataPt.Date < histItem.Date {
				// insert new entry
				newDataPt := SectorHistoryData{
					Date: histItem.Date,
				}
				appendCorrectSector(&newDataPt, sector.Id.Name, histItem.Total)
				*data = append((*data)[:i+1], (*data)[i:]...)
				(*data)[i] = newDataPt
				found = true
				break
			}
		}

		if !found {
			// insert new entry
			newDataPt := SectorHistoryData{
				Date: histItem.Date,
			}
			appendCorrectSector(&newDataPt, sector.Id.Name, histItem.Total)
			*data = append(*data, newDataPt)
		}
	}
}

func appendCorrectSector(pt *SectorHistoryData, sector string, value float64) {
	if sector == "Real Estate" {
		pt.RealEstate += value
	} else if sector == "Healthcare" {
		pt.Healthcare += value
	} else if sector == "Consumer Defensive" {
		pt.ConsumerDefensive += value
	} else if sector == "Fund" {
		pt.Fund += value
	} else if sector == "Energy" {
		pt.Energy += value
	} else if sector == "Basic Materials" {
		pt.BasicMaterials += value
	} else if sector == "Industrials" {
		pt.Industrials += value
	} else if sector == "Technology" {
		pt.Technology += value
	} else if sector == "Utilities" {
		pt.Utilities += value
	} else if sector == "Consumer Cyclical" {
		pt.ConsumerCyclical += value
	} else if sector == "Communication Services" {
		pt.CommunicationServices += value
	} else if sector == "Financial Services" {
		pt.FinancialServices += value
	}
}

type SectorHistoryData struct {
	Date                  string  `json:"date"`
	RealEstate            float64 `json:"realEstate"`
	Healthcare            float64 `json:"healthcare"`
	ConsumerDefensive     float64 `json:"consumerDefensive"`
	Fund                  float64 `json:"fund"`
	Energy                float64 `json:"energy"`
	BasicMaterials        float64 `json:"basicMaterials"`
	Industrials           float64 `json:"industrials"`
	Technology            float64 `json:"technology"`
	Utilities             float64 `json:"utilities"`
	ConsumerCyclical      float64 `json:"consumerCyclical"`
	CommunicationServices float64 `json:"communicationServices"`
	FinancialServices     float64 `json:"financialServices"`
}
