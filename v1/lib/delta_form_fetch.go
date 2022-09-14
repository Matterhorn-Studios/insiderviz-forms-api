package lib

import (
	"context"

	"github.com/Matterhorn-Studios/insiderviz-forms-api/v1/v1_database"
	"github.com/Matterhorn-Studios/insidervizforms/iv_models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DeltaFormFetch(filter primitive.D, opts *options.FindOptions) ([]iv_models.DB_DeltaForm, error) {
	deltaCollection := v1_database.GetCollection("DeltaForm")
	deltaForms := make([]iv_models.DB_DeltaForm, 0)
	data, errFetch := deltaCollection.Find(context.TODO(), filter, opts)
	if errFetch != nil {
		return nil, errFetch
	}
	errParse := data.All(context.TODO(), &deltaForms)
	if errParse != nil {
		return nil, errParse
	}
	return deltaForms, nil
}
