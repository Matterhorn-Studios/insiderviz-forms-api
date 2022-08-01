package main

import (
	"context"
	"gin/config"
	"gin/structs"
	"gin/utils"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupRouter() *gin.Engine {

	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.Use(authCheck())

	r.GET("/search", func(c *gin.Context) {
		// get the query
		query := c.Query("query")

		// setup the filter
		searchFilter := bson.D{{Key: "$search", Value: bson.D{
			{Key: "index", Value: "issuer_name"},
			{Key: "text", Value: bson.D{
				{Key: "query", Value: query},
				{Key: "path", Value: bson.A{"name", "cik", "ein", "tickers"}},
				{Key: "fuzzy", Value: bson.D{}},
			}},
		}}}

		limitFilter := bson.D{{Key: "$limit", Value: 10}}

		cursor, err := config.GetCollection(config.DB, "Issuer").Aggregate(context.TODO(), mongo.Pipeline{searchFilter, limitFilter})
		if err != nil {
			panic(err)
		}

		var issuers []structs.DB_Issuer_Doc
		if err = cursor.All(context.TODO(), &issuers); err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, issuers)
	})

	r.GET("/dailySentiment", func(c *gin.Context) {

		// get the start date from query
		startDate := c.Query("startDate")

		// setup the aggregate pipeline
		matchStage := bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "periodOfReport", Value: bson.D{{Key: "$gte", Value: startDate}}},
			}},
		}

		projectStage := bson.D{
			{Key: "$project", Value: bson.D{
				{
					Key: "SumSells", Value: bson.D{
						{Key: "$cond", Value: bson.A{
							bson.D{{
								Key: "$eq", Value: bson.A{"$buyOrSell", "Sell"},
							}},
							"$netTotal", 0,
						}},
					},
				},
				{
					Key: "SumBuys", Value: bson.D{
						{Key: "$cond", Value: bson.A{
							bson.D{{
								Key: "$eq", Value: bson.A{"$buyOrSell", "Buy"},
							}},
							"$netTotal", 0,
						}},
					},
				},
				{
					Key: "CountSells", Value: bson.D{
						{Key: "$cond", Value: bson.A{
							bson.D{{
								Key: "$eq", Value: bson.A{"$buyOrSell", "Sell"},
							}},
							1, 0,
						}},
					},
				},
				{
					Key: "CountBuys", Value: bson.D{
						{Key: "$cond", Value: bson.A{
							bson.D{{
								Key: "$eq", Value: bson.A{"$buyOrSell", "Buy"},
							}},
							1, 0,
						}},
					},
				},
				{
					Key: "periodOfReport", Value: 1,
				},
			}},
		}

		groupStage := bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: "$periodOfReport"},
				{Key: "sellAmount", Value: bson.D{
					{Key: "$sum", Value: "$SumSells"},
				}},
				{Key: "buyAmount", Value: bson.D{
					{Key: "$sum", Value: "$SumBuys"},
				}},
				{Key: "sellCount", Value: bson.D{
					{Key: "$sum", Value: "$CountSells"},
				}},
				{Key: "buyCount", Value: bson.D{
					{Key: "$sum", Value: "$CountBuys"},
				}},
			}},
		}

		orderState := bson.D{
			{Key: "$sort", Value: bson.D{{Key: "_id", Value: -1}}},
		}

		// run the aggregate
		cursor, err := config.GetCollection(config.DB, "DeltaForm").Aggregate(context.TODO(), mongo.Pipeline{matchStage, projectStage, groupStage, orderState})
		if err != nil {
			panic(err)
		}

		// display the results
		var results []bson.M
		if err = cursor.All(context.TODO(), &results); err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, results)
	})

	// get the top 50 DeltaForms from the last month
	r.GET("/delta", func(c *gin.Context) {
		// get the params
		buySellOrBoth := c.Query("buySellOrBoth")
		insiderCongressOrBoth := c.Query("insiderCongressOrBoth")

		var filterOne bson.E = bson.E{}
		if buySellOrBoth == "Buy" {
			filterOne = bson.E{Key: "buyOrSell", Value: "Buy"}
		} else if buySellOrBoth == "Sell" {
			filterOne = bson.E{Key: "buyOrSell", Value: "Sell"}
		}

		var filterTwo bson.E = bson.E{}
		if insiderCongressOrBoth == "Insider" {
			filterTwo = bson.E{Key: "formClass", Value: "Insider"}
		} else if insiderCongressOrBoth == "Congress" {
			filterTwo = bson.E{Key: "formClass", Value: "Congress"}
		}

		// offset and limit
		offset := c.Query("offset")
		limit := c.Query("limit")

		// parse the offset and limit
		var offsetInt int64 = 0
		var limitInt int64 = 50

		if offset != "" {
			offsetInt, _ = strconv.ParseInt(offset, 10, 64)
		}

		if limit != "" {
			limitInt, _ = strconv.ParseInt(limit, 10, 64)
		}

		filter := bson.D{{Key: "periodOfReport", Value: bson.D{{Key: "$gte", Value: "2022-06-15"}}}, filterOne, filterTwo}
		opts := options.Find().SetLimit(limitInt).SetSkip(offsetInt).SetSort(bson.D{{Key: "netTotal", Value: -1}})

		deltaForms, err := utils.DeltaFormFetch(filter, opts)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error1": err.Error()})
			return
		}

		c.JSON(http.StatusOK, deltaForms)
	})

	// get the 50 most recent forms
	r.GET("/recent", func(c *gin.Context) {
		// get the params
		buySellOrBoth := c.Query("buySellOrBoth")
		insiderCongressOrBoth := c.Query("insiderCongressOrBoth")

		var filterOne bson.E = bson.E{}
		if buySellOrBoth == "Buy" {
			filterOne = bson.E{Key: "buyOrSell", Value: "Buy"}
		} else if buySellOrBoth == "Sell" {
			filterOne = bson.E{Key: "buyOrSell", Value: "Sell"}
		}

		var filterTwo bson.E = bson.E{}
		if insiderCongressOrBoth == "Insider" {
			filterTwo = bson.E{Key: "formClass", Value: "Insider"}
		} else if insiderCongressOrBoth == "Congress" {
			filterTwo = bson.E{Key: "formClass", Value: "Congress"}
		}

		// offset and limit
		offset := c.Query("offset")
		limit := c.Query("limit")

		// parse the offset and limit
		var offsetInt int64 = 0
		var limitInt int64 = 50

		if offset != "" {
			offsetInt, _ = strconv.ParseInt(offset, 10, 64)
		}

		if limit != "" {
			limitInt, _ = strconv.ParseInt(limit, 10, 64)
		}

		filter := bson.D{filterOne, filterTwo}
		opts := options.Find().SetLimit(limitInt).SetSkip(offsetInt).SetSort(bson.D{{Key: "periodOfReport", Value: -1}})

		deltaForms, err := utils.DeltaFormFetch(filter, opts)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error1": err.Error()})
			return
		}

		c.JSON(http.StatusOK, deltaForms)
	})

	// get all the forms from one reporter
	r.GET("/reporter/:cik", func(c *gin.Context) {
		cik := c.Param("cik")

		filter := bson.D{{Key: "reporters.reporterCik", Value: cik}}
		opts := options.Find().SetSort(bson.D{{Key: "periodOfReport", Value: -1}})

		deltaForms, err := utils.DeltaFormFetch(filter, opts)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error1": err.Error()})
			return
		}

		// get the reporter's information
		reporterCollection := config.GetCollection(config.DB, "Reporter")
		filter = bson.D{{Key: "cik", Value: cik}}

		issuerInfo := reporterCollection.FindOne(context.TODO(), filter)

		if issuerInfo.Err() != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error1": issuerInfo.Err().Error()})
			return
		}

		// unmarshal the reporter info
		var reporter bson.M
		if err = issuerInfo.Decode(&reporter); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error1": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"forms": deltaForms, "info": reporter})

	})

	// get all the forms from one issuer
	r.GET("/issuer/:cik", func(c *gin.Context) {
		cik := c.Param("cik")

		filter := bson.D{{Key: "issuer.issuerCik", Value: cik}}
		opts := options.Find().SetSort(bson.D{{Key: "periodOfReport", Value: -1}})

		deltaForms, err := utils.DeltaFormFetch(filter, opts)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error1": err.Error()})
			return
		}

		// get the issuer's information
		issuerCollection := config.GetCollection(config.DB, "Issuer")
		filter = bson.D{{Key: "cik", Value: cik}}

		issuerInfo := issuerCollection.FindOne(context.TODO(), filter)

		if issuerInfo.Err() != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error1": issuerInfo.Err().Error()})
			return
		}

		// unmarshal the issuer info
		var issuer bson.M
		if err = issuerInfo.Decode(&issuer); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error1": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"forms": deltaForms, "info": issuer})
	})

	// get the featured issuers
	r.GET("/featuredIssuers", func(c *gin.Context) {
		var featuredIssuers []structs.DB_FeaturedIssuer
		cursor, err := config.GetCollection(config.DB, "FeaturedIssuer").Find(context.TODO(), bson.D{{}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error1": err.Error()})
			return
		}

		errParse := cursor.All(context.TODO(), &featuredIssuers)
		if errParse != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error1": errParse.Error()})
			return
		}

		c.JSON(http.StatusOK, featuredIssuers)
	})

	return r
}

func authCheck() gin.HandlerFunc {
	return func(c *gin.Context) {

		// check the header
		sentAuth := c.GetHeader("Authorization")
		realAuth := os.Getenv("AUTH_TOKEN")

		if sentAuth != realAuth {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		c.Next()

	}
}

func main() {
	r := setupRouter()

	// turn on DB
	config.ConnectDB()

	port := os.Getenv("PORT")

	// check if port is empty
	if port == "" {
		port = "8080"
	}

	// run on port
	r.Run(":" + port)
}
