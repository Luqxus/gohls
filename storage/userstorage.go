package storage

import (
	"context"

	"github.com/luqus/livespace/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthenticationStore interface {
	CreateUser(ctx context.Context, user *types.User) error
	FetchUser(ctx context.Context, email string) (*types.User, error)
	FetchUserByUID(ctx context.Context, uid string) (*types.User, error)
	CheckEmailExists(ctx context.Context, email string) (bool, error)
	CheckUsernameExists(ctx context.Context, username string) (bool, error)
}

type MongoAuthenticationStore struct {
	collection *mongo.Collection
}

// func NewMongoAuthenticationStore() *MongoAuthenticationStore {
// 	return
// }

// TODO: check if email exists in database
func (auth *MongoAuthenticationStore) CheckEmailExists(ctx context.Context, email string) (bool, error) {

	filter := primitive.D{primitive.E{Key: "email", Value: email}}
	result, err := auth.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	if result <= 0 {
		// TODO: if email not exists return false
		return false, nil
	}

	// TODO: if user exists return true
	return true, nil
}

// TODO: check if username exists in database
func (auth *MongoAuthenticationStore) CheckUsernameExists(ctx context.Context, username string) (bool, error) {
	filter := primitive.D{primitive.E{Key: "username", Value: username}}
	result, err := auth.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	if result <= 0 {
		// TODO: if username not exists return false
		return false, nil
	}
	// TODO: if username exists return true
	return true, nil
}

// TODO: create || commit user on user registration
func (auth *MongoAuthenticationStore) CreateUser(ctx context.Context, user *types.User) error {
	_, err := auth.collection.InsertOne(ctx, user)
	return err
}

// TODO: fetch user data on user login
func (auth *MongoAuthenticationStore) FetchUser(ctx context.Context, email string) (*types.User, error) {
	filter := primitive.D{primitive.E{Key: "email", Value: email}}

	user := new(types.User)

	err := auth.collection.FindOne(ctx, filter).Decode(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (auth *MongoAuthenticationStore) FetchUserByUID(ctx context.Context, uid string) (*types.User, error) {
	filter := primitive.D{primitive.E{Key: "uid", Value: uid}}
	user := new(types.User)

	err := auth.collection.FindOne(ctx, filter).Decode(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
