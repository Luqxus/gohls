package storage

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func NewStore() (*MongoAuthenticationStore, *MongoMetaDataStore) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// dbUrl := os.Getenv("DBURL")
	dbUrl := "mongodb://localhost:27017/"

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbUrl))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	metaDataCollection := client.Database("livespace").Collection("MetaData")
	userCollection := client.Database("livespace").Collection("User")
	return &MongoAuthenticationStore{
			collection: userCollection,
		}, &MongoMetaDataStore{
			collection: metaDataCollection,
		}

}
