package main

import (
	"gin/config"
	"gin/utils"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupRouter() *gin.Engine {

	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// get the top 50 DeltaForms from the last month
	r.GET("/delta", func(c *gin.Context) {
		filter := bson.D{{Key: "periodOfReport", Value: bson.D{{Key: "$gte", Value: "2022-06-15"}}}}
		opts := options.Find().SetLimit(50).SetSort(bson.D{{Key: "netTotal", Value: -1}})

		deltaForms, err := utils.DeltaFormFetch(filter, opts)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error1": err.Error()})
			return
		}

		c.JSON(http.StatusOK, deltaForms)
	})

	// get the 50 most recent forms
	r.GET("/recent", func(c *gin.Context) {
		filter := bson.D{}
		opts := options.Find().SetLimit(25).SetSort(bson.D{{Key: "periodOfReport", Value: -1}})

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

		c.JSON(http.StatusOK, deltaForms)

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

		c.JSON(http.StatusOK, deltaForms)
	})

	return r
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
