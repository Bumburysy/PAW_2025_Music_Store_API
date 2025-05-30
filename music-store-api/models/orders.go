package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	OrderStatusPending    = "pending"
	OrderStatusProcessing = "processing"
	OrderStatusShipped    = "shipped"
	OrderStatusCompleted  = "completed"
	OrderStatusCancelled  = "cancelled"
)

// Order reprezentuje zamówienie użytkownika
// swagger:model Order
type Order struct {
	// ID zamówienia (unikalny identyfikator)
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	// ID użytkownika, który złożył zamówienie
	UserID primitive.ObjectID `bson:"user_id" json:"user_id"`
	// Lista pozycji w zamówieniu
	Items []OrderItem `bson:"items" json:"items"`
	// Całkowita wartość zamówienia
	Total float64 `bson:"total" json:"total"`
	// Status zamówienia (pending, processing, shipped, completed, cancelled)
	Status string `bson:"status" json:"status"`
	// Data utworzenia zamówienia
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	// Data ostatniej aktualizacji zamówienia
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
	// Dane do wysyłki
	Shipping ShippingDetails `bson:"shipping" json:"shipping"`
}

// OrderItem reprezentuje pojedynczą pozycję zamówienia
// swagger:model OrderItem
type OrderItem struct {
	// ID albumu w zamówieniu
	AlbumID primitive.ObjectID `bson:"album_id" json:"album_id"`
	// Ilość sztuk albumu
	Quantity int `bson:"quantity" json:"quantity"`
	// Cena jednostkowa albumu w momencie zamówienia
	Price float64 `bson:"price" json:"price"`
}
