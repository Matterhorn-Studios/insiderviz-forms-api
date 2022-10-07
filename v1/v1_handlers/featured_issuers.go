package v1_handlers

import (
	"context"

	"github.com/Matterhorn-Studios/insiderviz-forms-api/v1/models"
	"github.com/Matterhorn-Studios/insiderviz-forms-api/v1/v1_database"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func FeaturedIssuers(c *fiber.Ctx) error {
	var featuredIssuers []models.DB_FeaturedIssuer
	cursor, err := v1_database.GetCollection("FeaturedIssuer").Find(context.TODO(), bson.D{{}})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}

	errParse := cursor.All(context.TODO(), &featuredIssuers)
	if errParse != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}

	return c.JSON(featuredIssuers)
}
