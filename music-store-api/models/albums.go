package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Album reprezentuje album muzyczny w sklepie
// swagger:model Album
type Album struct {
	// ID albumu
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	// Tytuł albumu
	Title string `bson:"title" json:"title"`
	// Nazwa wykonawcy
	Artist string `bson:"artist" json:"artist"`
	// Gatunek muzyczny
	Genre string `bson:"genre" json:"genre"`
	// Opis albumu
	Description string `bson:"description,omitempty" json:"description,omitempty"`
	// Data wydania
	ReleaseDate time.Time `bson:"release_date" json:"release_date"`
	// Lista utworów
	Tracks []string `bson:"tracks,omitempty" json:"tracks,omitempty"`
	// Cena albumu
	Price float64 `bson:"price" json:"price"`
	// Ilość dostępnych sztuk
	Quantity int `bson:"quantity" json:"quantity"`
	// Data utworzenia wpisu
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	// Data ostatniej aktualizacji
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
	// URL do okładki albumu
	CoverURL string `bson:"cover_url,omitempty" json:"cover_url,omitempty"`
}
