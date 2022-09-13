package aggregation

import (
	"net/http"
	"sync"

	"github.com/Matterhorn-Studios/insiderviz-forms-api/v1/utils"
	"github.com/Matterhorn-Studios/insidervizforms/iv_models"
	"github.com/gin-gonic/gin"
)

type SafeSectorList struct {
	mu   sync.Mutex
	data []iv_models.Top_From_Sector
	err  error
}

func (s *SafeSectorList) GetDataForSector(startDate string, sector string, wg *sync.WaitGroup) {

	data, err := utils.TopFromSector(startDate, sector)
	if err != nil {
		s.err = err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	defer wg.Done()

	s.data = append(s.data, data)

}

func TopFromSector(c *gin.Context) {
	// get the start date from the query
	startDate := c.Query("startDate")

	sectors := []string{"Real Estate", "Healthcare", "Consumer Defensive", "Fund", "Energy", "Basic Materials", "Industrials", "Technology", "Utilities", "Consumer Cyclical", "Communication Services", "Financial Services"}

	s := SafeSectorList{data: make([]iv_models.Top_From_Sector, 0), err: nil}

	var wg sync.WaitGroup

	for _, sector := range sectors {
		wg.Add(1)
		go s.GetDataForSector(startDate, sector, &wg)
	}

	wg.Wait()

	if s.err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": s.err.Error()})
		return
	}

	c.JSON(http.StatusOK, s.data)
}
