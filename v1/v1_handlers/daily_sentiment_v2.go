package v1_handlers

import (
	"context"

	"github.com/Matterhorn-Studios/insiderviz-forms-api/v1/models"
	"github.com/Matterhorn-Studios/insiderviz-forms-api/v1/v1_database"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func DailySentimentV2(c *fiber.Ctx) error {
	var dailySentiment []models.DB_SentimentDay
	cursor, err := v1_database.GetCollection("DailySentiment").Find(context.TODO(), bson.D{{}})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}

	errParse := cursor.All(context.TODO(), &dailySentiment)
	if errParse != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": errParse})
	}

	return c.JSON(dailySentiment)
}
