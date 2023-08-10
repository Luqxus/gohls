package storage

import (
	"context"

	"github.com/luqus/livespace/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MetaDataStore interface {
	CommitMetaData(ctx context.Context, metaData *types.MetaData) error
	FetchVideos(ctx context.Context) ([]*types.MetaData, error)
	FetchVideoMetaData(ctx context.Context, videoID string) (*types.MetaData, error)
}

type MongoMetaDataStore struct {
	collection *mongo.Collection
}

func (store *MongoMetaDataStore) CommitMetaData(ctx context.Context, metaData *types.MetaData) error {
	_, err := store.collection.InsertOne(ctx, metaData)

	return err
}

// TODO: fetch all videos in database
func (store *MongoMetaDataStore) FetchVideos(ctx context.Context) ([]*types.MetaData, error) {
	cursor, err := store.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	results := make([]*types.MetaData, 0)

	if err := cursor.All(ctx, results); err != nil {
		return nil, err
	}

	return results, nil

}

// TODO: fetch video meta data by video id
func (store *MongoMetaDataStore) FetchVideoMetaData(ctx context.Context, videoID string) (*types.MetaData, error) {
	filter := primitive.D{primitive.E{Key: "video_id", Value: videoID}}
	metaData := new(*types.MetaData)

	err := store.collection.FindOne(ctx, filter).Decode(metaData)
	if err != nil {
		return nil, err
	}

	return *metaData, nil
}
