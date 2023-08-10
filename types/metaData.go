package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MetaData struct {
	ID          primitive.ObjectID `bson:"_id"`
	VideoID     string             `json:"video_id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	// Creator      *User
	CreateAt     time.Time `json:"created_at"`
	VideoUrl     string    `json:"video_url"`
	ThumbnailUrl string    `json:"thumbnail_url"`
}
