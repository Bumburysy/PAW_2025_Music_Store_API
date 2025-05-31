package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User reprezentuje użytkownika systemu
// swagger:model User
type User struct {
	// ID użytkownika (unikalny identyfikator)
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	// Imię użytkownika
	FirstName string `bson:"first_name" json:"first_name"`
	// Nazwisko użytkownika
	LastName string `bson:"last_name" json:"last_name"`
	// Adres email użytkownika
	Email string `bson:"email" json:"email"`
	// Numer telefonu
	PhoneNumber string `bson:"phone_number,omitempty"`
	// Hasło
	Password string `bson:"-" json:"password"`
	// Hash hasła (niewidoczny w API)
	PasswordHash string `bson:"password_hash" json:"-"`
	// Rola użytkownika (np. "employee", "customer", "admin")
	Role string `bson:"role" json:"role"`
	// Data utworzenia użytkownika
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	// Data ostatniej aktualizacji
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
	// Czy konto jest aktywne
	IsActive bool `bson:"is_active" json:"is_active"`
	// Dane adresowe
	ShippingDetails *ShippingDetails `bson:"shipping_details,omitempty" json:"shipping_details,omitempty"`
}
