package utils

import "github.com/Matterhorn-Studios/insidervizforms/iv_models"

func TopFromSector(startDate string) (iv_models.Top_From_Sector, error) {
	var result iv_models.Top_From_Sector
	return result, nil
}

type Top_From_Sector struct {
	Sector    string                  `json:"sector" bson:"sector"`
	Companies []Top_From_Sector_Entry `json:"companies" bson:"companies"`
}

type Top_From_Sector_Entry struct {
	Ticker      string `json:"ticker" bson:"ticker"`
	Name        string `json:"name" bson:"name"`
	Industry    string `json:"industry" bson:"industry"`
	TradeVolume int    `json:"tradeVolume" bson:"tradeVolume"`
}
