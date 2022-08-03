package utils

import (
	"net/http"

	"github.com/Matterhorn-Studios/insiderviz-forms-api/config"
	"github.com/gin-gonic/gin"
)

func AuthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {

		// check the header
		sentAuth := c.GetHeader("Authorization")
		realAuth := config.GetEnvVariable("AUTH_TOKEN")

		if sentAuth != realAuth {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		c.Next()

	}
}
