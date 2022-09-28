package v1

import (
	"github.com/Matterhorn-Studios/insiderviz-forms-api/v1/v1_database"
	"github.com/Matterhorn-Studios/insiderviz-forms-api/v1/v1_handlers"
	"github.com/Matterhorn-Studios/insiderviz-forms-api/v1/v1_middleware"
	"github.com/gofiber/fiber/v2"
)

func AddV1Group(app *fiber.App) error {
	// start the db
	if err := v1_database.InitDb(); err != nil {
		return err
	}

	// create the group
	v1 := app.Group("/v1")

	// un-auth check
	v1.Get("/ping", v1_handlers.Ping)

	// auth check
	v1.Use(v1_middleware.AuthCheck)
	v1.Get("/pingAuth", v1_handlers.Ping)

	// aggregate
	a_g := v1.Group("/aggregate")
	{
		a_g.Get("/topFromSector", v1_handlers.TopFromSector)
		a_g.Get("/dailySentimentV2", v1_handlers.DailySentimentV2)
		a_g.Get("/top", v1_handlers.TopCompanies)
		a_g.Get("/featuredIssuers", v1_handlers.FeaturedIssuers)
		a_g.Get("/sectorHistory", v1_handlers.SectorHistory)
	}

	// search
	v1.Get("/search", v1_handlers.Search)

	// delta
	d_g := v1.Group("/delta")
	{
		d_g.Get("/topThisMonth", v1_handlers.TopThisMonth)
		d_g.Get("/recent", v1_handlers.Recent)
		d_g.Get("/deepFilter", v1_handlers.DeepFilter)
	}

	// single
	s_g := v1.Group("/single")
	{
		s_g.Get("/issuer/random", v1_handlers.RandomIssuer)
		s_g.Get("/issuer/:cik", v1_handlers.Issuer)
		s_g.Get("/issuer/graph/:cik", v1_handlers.IssuerGraph)
		s_g.Get("/reporter/random", v1_handlers.RandomReporter)
		s_g.Get("/reporter/:cik", v1_handlers.Reporter)
		s_g.Get("/reporter/holdings/:cik", v1_handlers.LatestThirteenF)
	}

	// email
	e_g := v1.Group("/email")
	{
		e_g.Get("/issuer/:cik/:date", v1_handlers.EmailFormIssuer)
		e_g.Get("/reporter/:cik/:date", v1_handlers.EmailFormReporter)
	}

	// csv
	c_g := v1.Group("/csv")
	{
		c_g.Get("/delta/issuer/:cik", v1_handlers.DeltaFormCsvIssuer)
		c_g.Get("/delta/reporter/:cik", v1_handlers.DeltaFormCsvReporter)
	}

	// config
	f_g := v1.Group("/config")
	{
		f_g.Get("/search/issuer", v1_handlers.SearchConfigIssuer)
	}

	return nil
}
