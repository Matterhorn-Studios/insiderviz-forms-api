package structs

type DB_Issuer_Doc struct {
	Name                 string      `json:"name" bson:"name"`
	Cik                  string      `json:"cik" bson:"cik"`
	Sic                  string      `json:"sic" bson:"sic"`
	SicDescription       string      `json:"secDescription" bson:"secDescription"`
	Ein                  string      `json:"ein" bson:"ein"`
	Tickers              []string    `json:"tickers" bson:"tickers"`
	Exchanges            []string    `json:"exchanges" bson:"exchanges"`
	FiscalYearEnd        string      `json:"fiscalYearEnd" bson:"fiscalYearEnd"`
	StateOfIncorporation string      `json:"stateOfIncorporation" bson:"stateOfIncorporation"`
	Phone                string      `json:"phone" bson:"phone"`
	StockData            []StockData `json:"stockData" bson:"stockData"`
}

// D = date, C = close, V = volume
type StockData struct {
	D string  `json:"d" bson:"d"`
	C float64 `json:"c" bson:"c"`
	V int     `json:"v" bson:"v"`
}
