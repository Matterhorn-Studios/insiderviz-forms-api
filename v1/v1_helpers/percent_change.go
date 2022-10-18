package v1_helpers

import (
	"github.com/Matterhorn-Studios/insiderviz-forms-api/v1/models"
)

func PercentChange(form models.DB_DeltaForm) int {
	// check if it was a buy or a sell
	if form.BuyOrSell == "Buy" {
		startShares := form.PostTransactionShares - form.SharesTraded

		if startShares == 0 {
			return 100
		}

		// calculate the percent change
		percentChange := (form.SharesTraded) / startShares

		return int(percentChange * 100)
	} else {
		startShares := form.PostTransactionShares + form.SharesTraded

		// calculate the percent change
		percentChange := (form.SharesTraded) / startShares

		return int(percentChange * 100)
	}
}
