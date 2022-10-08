package models

type DB_DeltaRuleTree struct {
	AccessionNumber  string `json:"accessionNumber" bson:"accessionNumber"`
	V_Footnote       bool   `json:"vFootnote" bson:"vFootnote"`
	V_HasFootnote    bool   `json:"vHasFootnote" bson:"vHasFootnote"`
	V_NetTotal       bool   `json:"vNetTotal" bson:"vNetTotal"`
	V_Price          bool   `json:"vPrice" bson:"vPrice"`
	V_Transactions   bool   `json:"vTransactions" bson:"vTransactions"`
	V_IssuerReporter bool   `json:"vIssuerReporter" bson:"vIssuerReporter"`
}
