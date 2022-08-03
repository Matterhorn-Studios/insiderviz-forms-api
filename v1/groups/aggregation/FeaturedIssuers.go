package aggregation

import (
	"context"
	"net/http"

	"github.com/Matterhorn-Studios/insiderviz-forms-api/config"
	"github.com/Matterhorn-Studios/insiderviz-forms-api/v1/structs"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func FeaturedIssuers(c *gin.Context) {
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
}
