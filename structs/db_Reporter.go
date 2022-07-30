package structs

type DB_Reporter_Doc struct {
	Name          string `json:"name" bson:"name"`
	Cik           string `json:"cik" bson:"cik"`
	IsCongressman bool   `json:"isCongressman" bson:"isCongressman"`
	Party         string `json:"party" bson:"party"`
}
