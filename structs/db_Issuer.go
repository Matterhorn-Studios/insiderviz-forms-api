package structs

type DB_Issuer_Doc struct {
	Name                 string   `json:"name" bson:"name"`
	Cik                  string   `json:"cik" bson:"cik"`
	Sic                  string   `json:"sic" bson:"sic"`
	SicDescription       string   `json:"secDescription" bson:"secDescription"`
	Ein                  string   `json:"ein" bson:"ein"`
	Tickers              []string `json:"tickers" bson:"tickers"`
	Exchanges            []string `json:"exchanges" bson:"exchanges"`
	FiscalYearEnd        string   `json:"fiscalYearEnd" bson:"fiscalYearEnd"`
	StateOfIncorporation string   `json:"stateOfIncorporation" bson:"stateOfIncorporation"`
	Phone                string   `json:"phone" bson:"phone"`
}
