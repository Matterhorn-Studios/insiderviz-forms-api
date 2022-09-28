package v1_database

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database
var client *mongo.Client

func GetCollection(name string) *mongo.Collection {
	return db.Collection(name)
}

func GetNewSession() (mongo.Session, error) {
	return client.StartSession()
}

func InitDb() error {
	// generate client
	var err error
	client, err = mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		return err
	}

	// connect
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return err
	}

	//ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	fmt.Println("Connected to MongoDB")

	db = client.Database("insiderviz")

	return nil
}
