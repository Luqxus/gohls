package storage

import "go.mongodb.org/mongo-driver/bson/primitive"

type UniqueID struct {
	id primitive.ObjectID
}

func NewID() *UniqueID {
	return &UniqueID{
		id: primitive.NewObjectID(),
	}
}

func (u *UniqueID) String() string {
	return u.id.Hex()
}

func (u *UniqueID) ID() primitive.ObjectID {
	return u.id
}
