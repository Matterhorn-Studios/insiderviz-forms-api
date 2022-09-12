package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Matterhorn-Studios/insiderviz-forms-api/database"
	v1 "github.com/Matterhorn-Studios/insiderviz-forms-api/v1"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("STARTING SERVER...")

	err := run()

	if err != nil {
		fmt.Println("SERVER ERROR:", err)
		os.Exit(1)
	}
}

func run() error {
	// ENV SETUP
	if err := initEnv(); err != nil {
		return err
	}

	// DB SETUP
	if err := database.InitDb(); err != nil {
		return err
	}

	// ROUTER SETUP
	r := setupRouter()

	// PORT SETUP
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// RUN ON PORT
	r.Run(":" + port)

	return nil
}

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

func initEnv() error {
	if os.Getenv("GO_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			return err
		}
	}

	return nil
}
