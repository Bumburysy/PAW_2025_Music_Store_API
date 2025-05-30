package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Cart reprezentuje koszyk użytkownika
// swagger:model Cart
type Cart struct {
	// ID koszyka (unikalny identyfikator)
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	// ID użytkownika, do którego należy koszyk
	UserID primitive.ObjectID `bson:"user_id" json:"user_id"`
	// Lista pozycji w koszyku
	Items []CartItem `bson:"items" json:"items"`
	// Sumaryczna wartości koszyka
	Total float64 `bson:"total,omitempty" json:"total,omitempty"`
}

// CartItem reprezentuje pojedynczą pozycję w koszyku
// swagger:model CartItem
type CartItem struct {
	// ID albumu dodanego do koszyka
	AlbumID primitive.ObjectID `bson:"album_id" json:"album_id"`
	// Ilość sztuk albumu w koszyku
	Quantity int `bson:"quantity" json:"quantity"`
}
