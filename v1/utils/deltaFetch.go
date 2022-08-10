package utils

import (
	"context"

	"github.com/Matterhorn-Studios/insiderviz-forms-api/config"
	"github.com/Matterhorn-Studios/insiderviz-forms-api/v1/structs"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var deltaCollection *mongo.Collection = config.GetCollection(config.DB, "DeltaForm")

func DeltaFormFetch(filter primitive.D, opts *options.FindOptions) ([]structs.DB_DeltaForm, error) {
	deltaForms := make([]structs.DB_DeltaForm, 0)
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
