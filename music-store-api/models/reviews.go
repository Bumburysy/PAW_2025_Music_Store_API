package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Review reprezentuje recenzję albumu
// swagger:model Review
type Review struct {
	// ID recenzji (unikalny identyfikator)
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	// ID albumu, którego dotyczy recenzja
	AlbumID primitive.ObjectID `bson:"album_id" json:"album_id"`
	// ID użytkownika, który dodał recenzję
	UserID primitive.ObjectID `bson:"user_id" json:"user_id"`
	// Ocena albumu (np. od 1 do 5)
	Rating int `bson:"rating" json:"rating"`
	// Komentarz do recenzji
	Comment string `bson:"comment" json:"comment"`
	// Data utworzenia recenzji
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}
