package top

import (
	"context"

	iv_structs "github.com/Matterhorn-Studios/insiderviz-backend_structs"
	"github.com/Matterhorn-Studios/insiderviz-forms-api/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TopInstitutions(c *gin.Context) {
	reporterCollection := database.GetCollection("Reporter")

	// get the top institutions, ordered by last13FTotal descending
	opts := options.Find().SetSort(bson.D{{Key: "last13FTotal", Value: -1}}).SetLimit(10)

	cur, err := reporterCollection.Find(context.TODO(), bson.D{{}}, opts)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var topInstitutions []iv_structs.DB_Reporter_Doc

	cur.All(context.TODO(), &topInstitutions)

	c.JSON(200, topInstitutions)
}
