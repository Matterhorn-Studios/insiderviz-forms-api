package v1

import (
	"github.com/Matterhorn-Studios/insiderviz-forms-api/v1/groups/aggregation"
	"github.com/Matterhorn-Studios/insiderviz-forms-api/v1/groups/delta"
	"github.com/Matterhorn-Studios/insiderviz-forms-api/v1/groups/health"
	"github.com/Matterhorn-Studios/insiderviz-forms-api/v1/groups/search"
	"github.com/Matterhorn-Studios/insiderviz-forms-api/v1/groups/single"
	"github.com/Matterhorn-Studios/insiderviz-forms-api/v1/utils"
	"github.com/gin-gonic/gin"
)

func AddGroup(r *gin.Engine) {
	v1 := r.Group("/v1")
	{
		// no auth check
		v1.GET("/ping", health.Ping)

		// auth check
		v1.Use(utils.AuthCheck())
		v1.GET("/pingAuth", health.Ping)

		// delta
		d_g := v1.Group("/delta")
		{
			d_g.GET("/topThisMonth", delta.TopThisMonth)
			d_g.GET("/recent", delta.Recent)
			d_g.GET("/deepFilter", delta.DeepFilter)
		}

		// single
		s_g := v1.Group("/single")
		{
			s_g.GET("/issuer/:cik", single.Issuer)
			s_g.GET("/issuer/random", single.RandomIssuer)
			s_g.GET("/issuer/graph/:cik", single.IssuerGraph)
			s_g.GET("/reporter/:cik", single.Reporter)
			s_g.GET("/reporter/random", single.RandomReporter)
		}

		// search
		se_g := v1.Group("/search")
		{
			se_g.GET("/", search.Search)
		}

		// aggregate
		a_g := v1.Group("/aggregate")
		{
			a_g.GET("/top", aggregation.Top)
			a_g.GET("/featuredIssuers", aggregation.FeaturedIssuers)
			a_g.GET("/dailySentiment", aggregation.DailySentiment)
		}
	}
}
