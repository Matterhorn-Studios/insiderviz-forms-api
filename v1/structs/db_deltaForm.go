package structs

type DB_DeltaForm struct {
	AccessionNumber       string        `json:"accessionNumber" bson:"accessionNumber"`
	FormClass             string        `json:"formClass" bson:"formClass"`
	PeriodOfReport        string        `json:"periodOfReport" bson:"periodOfReport"`
	AveragePricePerShare  float32       `json:"averagePricePerShare" bson:"averagePricePerShare"`
	NetTotal              float32       `json:"netTotal" bson:"netTotal"`
	SharesTraded          float32       `json:"sharesTraded" bson:"sharesTraded"`
	PostTransactionShares float32       `json:"postTransactionShares" bson:"postTransactionShares"`
	BuyOrSell             string        `json:"buyOrSell" bson:"buyOrSell"`
	Url                   string        `json:"url" bson:"url"`
	Issuer                DB_Issuer     `json:"issuer" bson:"issuer"`
	Reporters             []DB_Reporter `json:"reporters" bson:"reporters"`
}
