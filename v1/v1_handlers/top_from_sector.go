package v1_handlers

import (
	"sync"

	"github.com/Matterhorn-Studios/insiderviz-forms-api/v1/lib"
	"github.com/Matterhorn-Studios/insiderviz-forms-api/v1/models"
	"github.com/gofiber/fiber/v2"
)

type SafeSectorList struct {
	mu   sync.Mutex
	data []models.Top_From_Sector
	err  error
}

func (s *SafeSectorList) GetDataForSector(startDate string, sector string, wg *sync.WaitGroup) {

	data, err := lib.TopFromSector(startDate, sector)
	if err != nil {
		s.err = err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	defer wg.Done()

	s.data = append(s.data, data)

}

func TopFromSector(c *fiber.Ctx) error {
	// get the start date from the query
	startDate := c.Query("startDate")

	sectors := []string{"Real Estate", "Healthcare", "Consumer Defensive", "Fund", "Energy", "Basic Materials", "Industrials", "Technology", "Utilities", "Consumer Cyclical", "Communication Services", "Financial Services"}

	s := SafeSectorList{data: make([]models.Top_From_Sector, 0), err: nil}

	var wg sync.WaitGroup

	for _, sector := range sectors {
		wg.Add(1)
		go s.GetDataForSector(startDate, sector, &wg)
	}

	wg.Wait()

	if s.err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": s.err})
	}

	return c.JSON(s.data)
}
