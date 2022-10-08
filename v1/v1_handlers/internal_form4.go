package v1_handlers

import (
	"context"
	"sync"

	"github.com/Matterhorn-Studios/insiderviz-forms-api/v1/models"
	"github.com/Matterhorn-Studios/insiderviz-forms-api/v1/v1_database"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func InternalForm4(c *fiber.Ctx) error {
	// get the accession number from the query
	accessionNumber := c.Query("accessionNumber")

	// error handling
	errC := make(chan error)
	var wg sync.WaitGroup
	wg.Add(3)

	// fetch the raw, delta, and tree
	var raw models.DB_BaseForm4
	go func() {
		defer wg.Done()
		// get DB session
		session, err := v1_database.GetNewSession()
		if err != nil {
			errC <- err
			return
		}

		// get the raw form
		rawCollection := session.Client().Database("insiderviz").Collection("RawForm4")
		filter := bson.D{{Key: "accessionNumber", Value: accessionNumber}}
		doc := rawCollection.FindOne(context.Background(), filter)
		if doc.Err() != nil {
			errC <- doc.Err()
			return
		}

		// decode the raw form
		err = doc.Decode(&raw)
		if err != nil {
			errC <- err
			return
		}

	}()

	var delta models.DB_DeltaForm
	go func() {
		defer wg.Done()
		session, err := v1_database.GetNewSession()
		if err != nil {
			errC <- err
			return
		}

		// get the delta form
		deltaCollection := session.Client().Database("insiderviz").Collection("DeltaForm")
		filter := bson.D{{Key: "accessionNumber", Value: accessionNumber}}
		doc := deltaCollection.FindOne(context.Background(), filter)
		if doc.Err() != nil {
			return
		}

		// decode the delta form
		err = doc.Decode(&delta)
		if err != nil {
			errC <- err
			return
		}

	}()

	var tree models.DB_DeltaRuleTree
	go func() {
		defer wg.Done()
		session, err := v1_database.GetNewSession()
		if err != nil {
			errC <- err
			return
		}

		// get the delta form
		treeCollection := session.Client().Database("insiderviz").Collection("DeltaFormTree")
		filter := bson.D{{Key: "accessionNumber", Value: accessionNumber}}
		doc := treeCollection.FindOne(context.Background(), filter)
		if doc.Err() != nil {
			return
		}

		// decode the delta form
		err = doc.Decode(&tree)
		if err != nil {
			errC <- err
			return
		}

	}()

	// wait for all the goroutines to finish
	go func() {
		wg.Wait()
		close(errC)
	}()

	// handle sending the error
	for e := range errC {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": e.Error()})
	}

	// return the data
	if raw.Parsed {
		return c.JSON(fiber.Map{
			"raw":   raw,
			"delta": delta,
			"tree":  tree,
		})
	} else {
		return c.JSON(fiber.Map{
			"raw":   raw,
			"delta": nil,
			"tree":  tree,
		})

	}

}
