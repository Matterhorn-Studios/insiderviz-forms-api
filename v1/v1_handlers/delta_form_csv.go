package v1_handlers

import (
	"github.com/Matterhorn-Studios/insiderviz-forms-api/v1/v1_database"
	"github.com/Matterhorn-Studios/insidervizforms/iv_models"
	"github.com/gocarina/gocsv"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DeltaCsv struct {
	Date                  string  `csv:"date"`
	FormClass             string  `csv:"form_class"`
	Price                 float64 `csv:"price"`
	Shares                float64 `csv:"shares"`
	Total                 float64 `csv:"total"`
	PostTransactionShares float64 `csv:"post_transaction_shares"`
	IssuerName            string  `csv:"issuer_name"`
	IssuerTicker          string  `csv:"issuer_ticker"`
	ReporterName          string  `csv:"reporter_name"`
	ReporterTitle         string  `csv:"reporter_title"`
	Link                  string  `csv:"link"`
}

func DeltaFormCsvReporter(c *fiber.Ctx) error {
	cik := c.Params("cik")

	// fetch all delta forms for the provided issuer
	filter := bson.D{{Key: "reporters.reporterCik", Value: cik}}
	opts := options.Find().SetSort(bson.D{{Key: "periodOfReport", Value: -1}})
	cur, err := v1_database.GetCollection("DeltaForm").Find(c.Context(), filter, opts)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err})
	}

	var formList []iv_models.DB_DeltaForm
	if err := cur.All(c.Context(), &formList); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err})
	}

	csvList := make([]DeltaCsv, len(formList))

	for i, form := range formList {
		csvList[i] = DeltaCsv{
			FormClass:             form.FormClass,
			Date:                  form.PeriodOfReport,
			Price:                 float64(form.AveragePricePerShare),
			Shares:                float64(form.SharesTraded),
			Total:                 float64(form.NetTotal),
			PostTransactionShares: float64(form.PostTransactionShares),
			IssuerName:            form.Issuer.IssuerName,
			IssuerTicker:          form.Issuer.IssuerTicker,
			ReporterName:          form.Reporters[0].ReporterName,
			ReporterTitle:         form.Reporters[0].ReporterTitle,
			Link:                  form.Url,
		}

		if form.BuyOrSell == "Sell" {
			csvList[i].Total *= -1
		}
	}

	csvContent, err := gocsv.MarshalString(&csvList)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err})
	}

	return c.SendString(csvContent)
}

func DeltaFormCsvIssuer(c *fiber.Ctx) error {
	cik := c.Params("cik")

	// fetch all delta forms for the provided issuer
	filter := bson.D{{Key: "issuer.issuerCik", Value: cik}}
	opts := options.Find().SetSort(bson.D{{Key: "periodOfReport", Value: -1}})
	cur, err := v1_database.GetCollection("DeltaForm").Find(c.Context(), filter, opts)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err})
	}

	var formList []iv_models.DB_DeltaForm
	if err := cur.All(c.Context(), &formList); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err})
	}

	csvList := make([]DeltaCsv, len(formList))

	for i, form := range formList {
		csvList[i] = DeltaCsv{
			FormClass:             form.FormClass,
			Date:                  form.PeriodOfReport,
			Price:                 float64(form.AveragePricePerShare),
			Shares:                float64(form.SharesTraded),
			Total:                 float64(form.NetTotal),
			PostTransactionShares: float64(form.PostTransactionShares),
			IssuerName:            form.Issuer.IssuerName,
			IssuerTicker:          form.Issuer.IssuerTicker,
			ReporterName:          form.Reporters[0].ReporterName,
			ReporterTitle:         form.Reporters[0].ReporterTitle,
			Link:                  form.Url,
		}

		if form.BuyOrSell == "Sell" {
			csvList[i].Total *= -1
		}
	}

	csvContent, err := gocsv.MarshalString(&csvList)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err})
	}

	return c.SendString(csvContent)
}
