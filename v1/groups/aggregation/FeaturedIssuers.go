package aggregation

import (
	"context"
	"net/http"

	iv_structs "github.com/Matterhorn-Studios/insiderviz-backend_structs"
	"github.com/Matterhorn-Studios/insiderviz-forms-api/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func FeaturedIssuers(c *gin.Context) {
	var featuredIssuers []iv_structs.DB_FeaturedIssuer
	cursor, err := database.GetCollection("FeaturedIssuer").Find(context.TODO(), bson.D{{}})
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
