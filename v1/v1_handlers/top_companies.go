package v1_handlers

import (
	"sync"

	"github.com/Matterhorn-Studios/insiderviz-forms-api/v1/lib"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TopCompanies(c *fiber.Ctx) error {
	// get the start date from the query
	startDate := c.Query("startDate")

	var err error
	var mostBoughtInsider []primitive.M
	var mostSoldInsider []primitive.M
	var mostBoughtCongress []primitive.M
	var mostSoldCongress []primitive.M
	var wg sync.WaitGroup

	wg.Add(4)
	go func() {
		defer wg.Done()
		mostBoughtInsider, err = lib.MostBought(startDate, "Insider")
	}()

	go func() {
		defer wg.Done()
		mostSoldInsider, err = lib.MostSold(startDate, "Insider")
	}()

	go func() {
		defer wg.Done()
		mostBoughtCongress, err = lib.MostBought(startDate, "Congress")
	}()

	go func() {
		defer wg.Done()
		mostSoldCongress, err = lib.MostSold(startDate, "Congress")
	}()

	wg.Wait()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}

	returnObj := bson.M{
		"mostBoughtInsider":  mostBoughtInsider,
		"mostSoldInsider":    mostSoldInsider,
		"mostBoughtCongress": mostBoughtCongress,
		"mostSoldCongress":   mostSoldCongress,
	}

	return c.JSON(returnObj)
}
