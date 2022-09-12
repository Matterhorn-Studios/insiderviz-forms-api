package utils

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func AuthCheck() gin.HandlerFunc {
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
