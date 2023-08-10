package storage

import (
	"context"

	"github.com/luqus/livespace/types"
	"go.mongodb.org/mongo-driver/mongo"
)

type MetaDataStore interface {
	CommitMetaData(ctx context.Context, metaData *types.MetaData) error
}

type MongoMetaDataStore struct {
	collection *mongo.Collection
}

func (store *MongoMetaDataStore) CommitMetaData(ctx context.Context, metaData *types.MetaData) error {
	_, err := store.collection.InsertOne(ctx, metaData)

	return err
}
