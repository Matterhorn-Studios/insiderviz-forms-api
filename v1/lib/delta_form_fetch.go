package lib

import (
	"context"

	"github.com/Matterhorn-Studios/insiderviz-forms-api/v1/models"
	"github.com/Matterhorn-Studios/insiderviz-forms-api/v1/v1_database"
	"github.com/Matterhorn-Studios/insiderviz-forms-api/v1/v1_helpers"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DeltaFormFetch(filter primitive.D, opts *options.FindOptions) ([]models.DB_DeltaForm, error) {
	deltaCollection := v1_database.GetCollection("DeltaForm")
	deltaForms := make([]models.DB_DeltaForm, 0)
	data, errFetch := deltaCollection.Find(context.TODO(), filter, opts)
	if errFetch != nil {
		return nil, errFetch
	}
	errParse := data.All(context.TODO(), &deltaForms)
	if errParse != nil {
		return nil, errParse
	}
	for i := 0; i < len(deltaForms); i++ {
		deltaForms[i].PercentChange = v1_helpers.PercentChange(deltaForms[i])
	}

	return deltaForms, nil
}
