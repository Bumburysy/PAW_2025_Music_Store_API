package models

// ShippingDetails reprezentuje dane do wysyłki
// swagger:model ShippingDetails
type ShippingDetails struct {
	// Adres dostawy (np. ulica, nr domu/mieszkania)
	Address string `bson:"address" json:"address"`
	// Miasto dostawy
	City string `bson:"city" json:"city"`
	// Kod pocztowy dostawy
	PostalCode string `bson:"postal_code" json:"postal_code"`
	// Kraj dostawy
	Country string `bson:"country" json:"country"`
	// Numer telefonu kontaktowego
	PhoneNumber string `bson:"phone_number" json:"phone_number"`
}
