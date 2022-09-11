package utils

import (
	"context"

	iv_structs "github.com/Matterhorn-Studios/insiderviz-backend_structs"
	"github.com/Matterhorn-Studios/insiderviz-forms-api/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DeltaFormFetch(filter primitive.D, opts *options.FindOptions) ([]iv_structs.DB_DeltaForm, error) {
	var deltaCollection *mongo.Collection = database.GetCollection("DeltaForm")
	deltaForms := make([]iv_structs.DB_DeltaForm, 0)
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
