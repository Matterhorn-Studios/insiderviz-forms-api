package aggregation

import (
	"context"
	"net/http"

	"github.com/Matterhorn-Studios/insiderviz-forms-api/database"
	"github.com/Matterhorn-Studios/insidervizforms/iv_models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func DailySentimentV2(c *gin.Context) {
	var dailySentiment []iv_models.DB_SentimentDay
	cursor, err := database.GetCollection("DailySentiment").Find(context.TODO(), bson.D{{}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error1": err.Error()})
		return
	}

	errParse := cursor.All(context.TODO(), &dailySentiment)
	if errParse != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error1": errParse.Error()})
		return
	}

	c.JSON(http.StatusOK, dailySentiment)
}
