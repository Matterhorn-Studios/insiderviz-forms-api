package main

import (
	"net/http"
	"os"

	"github.com/Matterhorn-Studios/insiderviz-forms-api/config"
	v1 "github.com/Matterhorn-Studios/insiderviz-forms-api/v1"
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {

	r := gin.Default()

	// add the V1 group
	v1.AddGroup(r)

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
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
