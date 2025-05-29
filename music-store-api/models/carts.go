package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Cart struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	Items     []CartItem         `bson:"items" json:"items"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type CartItem struct {
	AlbumID  primitive.ObjectID `bson:"album_id" json:"album_id"`
	Quantity int                `bson:"quantity" json:"quantity"`
}
