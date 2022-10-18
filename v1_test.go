package main

import (
	"testing"

	"github.com/Matterhorn-Studios/insiderviz-forms-api/v1/models"
	"github.com/Matterhorn-Studios/insiderviz-forms-api/v1/v1_helpers"
)

func TestPercentChangeBuy(t *testing.T) {
	var testForm models.DB_DeltaForm
	testForm.BuyOrSell = "Buy"
	testForm.PostTransactionShares = 110
	testForm.SharesTraded = 10
	val := v1_helpers.PercentChange(testForm)

	if val != 10 {
		t.Error("Expected 10, got", val)
	}
}

func TestPercentChangeSell(t *testing.T) {
	var testForm models.DB_DeltaForm
	testForm.BuyOrSell = "Sell"
	testForm.PostTransactionShares = 90
	testForm.SharesTraded = 10
	val := v1_helpers.PercentChange(testForm)

	if val != 10 {
		t.Error("Expected 10, got", val)
	}
}

func TestPercentChangeBuyEmpty(t *testing.T) {
	var testForm models.DB_DeltaForm
	testForm.BuyOrSell = "Buy"
	testForm.PostTransactionShares = 100
	testForm.SharesTraded = 100
	val := v1_helpers.PercentChange(testForm)

	if val != 100 {
		t.Error("Expected 100, got", val)
	}
}

func TestPercentChangeSellEmpty(t *testing.T) {
	var testForm models.DB_DeltaForm
	testForm.BuyOrSell = "Sell"
	testForm.PostTransactionShares = 0
	testForm.SharesTraded = 100
	val := v1_helpers.PercentChange(testForm)

	if val != 100 {
		t.Error("Expected 100, got", val)
	}
}
