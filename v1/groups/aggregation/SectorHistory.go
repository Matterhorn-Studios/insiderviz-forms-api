package aggregation

import (
	"net/http"

	"github.com/Matterhorn-Studios/insiderviz-forms-api/v1/utils"
	"github.com/gin-gonic/gin"
)

func SectorHistory(c *gin.Context) {
	// get the start date from the query
	startDate := c.Query("startDate")

	data, err := utils.CalcSectorHistory(startDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(http.StatusOK, data)
}
