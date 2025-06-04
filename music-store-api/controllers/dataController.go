package controllers

import (
	"context"
	"encoding/json"
	"log"
	"music-store-api/config"
	"music-store-api/middleware"
	"music-store-api/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// LoadTestData ładuje dane testowe do bazy danych.
// @Summary Ładowanie danych testowych
// @Security BearerAuth
// @Description Wczytuje dane z plików JSON i wstawia je do kolekcji MongoDB.
// @Tags Data
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]string "Dane zostały wczytane"
// @Failure 500 {object} map[string]string "Błąd serwera lub pliku danych"
// @Router /data/load [post]
func LoadTestData(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := config.DB

	type LoadFile struct {
		FilePath   string
		Collection string
		Data       interface{}
		Process    func(interface{})
	}

	files := []LoadFile{
		{
			FilePath:   "data/albums.json",
			Collection: "albums",
			Data:       &[]models.Album{},
			Process: func(data interface{}) {
				albums := data.(*[]models.Album)
				now := time.Now()
				for i := range *albums {
					(*albums)[i].ID = primitive.NewObjectID()
					(*albums)[i].CreatedAt = now
					(*albums)[i].UpdatedAt = now
				}
			},
		},
		{
			FilePath:   "data/users.json",
			Collection: "users",
			Data:       &[]models.User{},
			Process: func(data interface{}) {
				users := data.(*[]models.User)
				now := time.Now()
				for i := range *users {
					(*users)[i].ID = primitive.NewObjectID()
					(*users)[i].CreatedAt = now
					(*users)[i].UpdatedAt = now

					hashed, err := middleware.HashPassword((*users)[i].Password)
					if err != nil {
						log.Printf("Błąd haszowania hasła użytkownika: %v", err)
						(*users)[i].PasswordHash = ""
					} else {
						(*users)[i].PasswordHash = hashed
					}

					(*users)[i].Password = ""
				}
			},
		},
	}

	for _, file := range files {
		byteValue, err := os.ReadFile(file.FilePath)
		if err != nil {
			log.Printf("Nie udało się otworzyć %s: %v", file.FilePath, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd odczytu pliku " + file.FilePath})
			return
		}

		err = json.Unmarshal(byteValue, file.Data)
		if err != nil {
			log.Printf("Błąd dekodowania %s: %v", file.FilePath, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd dekodowania pliku " + file.FilePath})
			return
		}

		file.Process(file.Data)

		if err := db.Collection(file.Collection).Drop(ctx); err != nil {
			log.Printf("Błąd przy czyszczeniu kolekcji %s: %v", file.Collection, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd czyszczenia kolekcji " + file.Collection})
			return
		}

		docs := toInterfaceSlice(file.Data)
		if len(docs) > 0 {
			_, err := db.Collection(file.Collection).InsertMany(ctx, docs)
			if err != nil {
				log.Printf("Błąd przy wstawianiu do %s: %v", file.Collection, err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd wstawiania do kolekcji " + file.Collection})
				return
			}
		}
	}

	var albums []models.Album
	var users []models.User

	cursor, err := db.Collection("albums").Find(ctx, bson.M{})
	if err != nil {
		log.Println("Błąd przy odczycie albumów:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd przy odczycie albumów"})
		return
	}
	if err := cursor.All(ctx, &albums); err != nil || len(albums) == 0 {
		log.Println("Błąd dekodowania albumów lub brak danych:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd przy dekodowaniu albumów"})
		return
	}

	cursor, err = db.Collection("users").Find(ctx, bson.M{})
	if err != nil {
		log.Println("Błąd przy odczycie użytkowników:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd przy odczycie użytkowników"})
		return
	}
	if err := cursor.All(ctx, &users); err != nil || len(users) == 0 {
		log.Println("Błąd dekodowania użytkowników lub brak danych:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd przy dekodowaniu użytkowników"})
		return
	}

	reviews := []models.Review{
		{
			ID:        primitive.NewObjectID(),
			AlbumID:   albums[0].ID,
			UserID:    users[0].ID,
			Rating:    5,
			Comment:   "Świetny album!",
			CreatedAt: time.Now(),
		},
		{
			ID:        primitive.NewObjectID(),
			AlbumID:   albums[1%len(albums)].ID,
			UserID:    users[1%len(users)].ID,
			Rating:    4,
			Comment:   "Fajny, ale mógłby być lepszy.",
			CreatedAt: time.Now(),
		},
		{
			ID:        primitive.NewObjectID(),
			AlbumID:   albums[2%len(albums)].ID,
			UserID:    users[2%len(users)].ID,
			Rating:    3,
			Comment:   "Nie do końca mój klimat.",
			CreatedAt: time.Now(),
		},
		{
			ID:        primitive.NewObjectID(),
			AlbumID:   albums[1].ID,
			UserID:    users[0].ID,
			Rating:    5,
			Comment:   "Kolejny hit! Polecam każdemu.",
			CreatedAt: time.Now(),
		},
	}

	orders := []models.Order{
		{
			ID:     primitive.NewObjectID(),
			UserID: users[0].ID,
			Items: []models.OrderItem{
				{
					AlbumID:  albums[0].ID,
					Quantity: 2,
					Price:    albums[0].Price,
				},
			},
			Total:     albums[0].Price * 2,
			Status:    "pending",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Shipping: models.ShippingDetails{
				Address:     "ul. Testowa 1",
				City:        "Warszawa",
				PostalCode:  "00-001",
				Country:     "Polska",
				PhoneNumber: "+48 600 000 001",
			},
		},
		{
			ID:     primitive.NewObjectID(),
			UserID: users[1].ID,
			Items: []models.OrderItem{
				{
					AlbumID:  albums[1].ID,
					Quantity: 1,
					Price:    albums[1].Price,
				},
				{
					AlbumID:  albums[2%len(albums)].ID,
					Quantity: 3,
					Price:    albums[2%len(albums)].Price,
				},
			},
			Total:     albums[1].Price*1 + albums[2%len(albums)].Price*3,
			Status:    "processing",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Shipping: models.ShippingDetails{
				Address:     "ul. Muzyczna 7",
				City:        "Kraków",
				PostalCode:  "30-002",
				Country:     "Polska",
				PhoneNumber: "+48 600 000 002",
			},
		},
		{
			ID:     primitive.NewObjectID(),
			UserID: users[2%len(users)].ID,
			Items: []models.OrderItem{
				{
					AlbumID:  albums[0].ID,
					Quantity: 1,
					Price:    albums[0].Price,
				},
			},
			Total:     albums[0].Price,
			Status:    "completed",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Shipping: models.ShippingDetails{
				Address:     "ul. Finalna 99",
				City:        "Gdańsk",
				PostalCode:  "80-003",
				Country:     "Polska",
				PhoneNumber: "+48 600 000 003",
			},
		},
	}

	if err := db.Collection("reviews").Drop(ctx); err != nil {
		log.Printf("Błąd przy czyszczeniu kolekcji reviews: %v", err)
	}
	if _, err := db.Collection("reviews").InsertMany(ctx, toInterfaceSlice(&reviews)); err != nil {
		log.Printf("Błąd przy wstawianiu recenzji: %v", err)
	}

	if err := db.Collection("orders").Drop(ctx); err != nil {
		log.Printf("Błąd przy czyszczeniu kolekcji orders: %v", err)
	}
	if _, err := db.Collection("orders").InsertMany(ctx, toInterfaceSlice(&orders)); err != nil {
		log.Printf("Błąd przy wstawianiu zamówień: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Dane zostały wczytane"})
}

func toInterfaceSlice(data interface{}) []interface{} {
	var res []interface{}
	switch v := data.(type) {
	case *[]models.Album:
		for _, item := range *v {
			res = append(res, item)
		}
	case *[]models.User:
		for _, item := range *v {
			res = append(res, item)
		}
	case *[]models.Review:
		for _, item := range *v {
			res = append(res, item)
		}
	case *[]models.Order:
		for _, item := range *v {
			res = append(res, item)
		}
	}
	return res
}
