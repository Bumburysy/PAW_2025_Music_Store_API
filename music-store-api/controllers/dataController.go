package controllers

import (
	"context"
	"encoding/json"
	"log"
	"music-store-api/config"
	"music-store-api/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func LoadTestData(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := config.DB

	type FileData struct {
		Filename   string
		Collection string
		Model      interface{}
	}

	files := []FileData{
		{"albums.json", "albums", &[]models.Album{}},
		{"users.json", "users", &[]models.User{}},
		{"orders.json", "orders", &[]models.Order{}},
		{"carts.json", "carts", &[]models.Cart{}},
		{"reviews.json", "reviews", &[]models.Review{}},
	}

	for _, file := range files {
		path := "data/" + file.Filename

		byteValue, err := os.ReadFile(path)
		if err != nil {
			log.Printf("Nie udało się otworzyć %s: %v", path, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd odczytu pliku " + file.Filename})
			return
		}

		err = json.Unmarshal(byteValue, file.Model)
		if err != nil {
			log.Printf("Błąd dekodowania %s: %v", path, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd dekodowania pliku " + file.Filename})
			return
		}

		setTimestamps(file.Model)

		if err := db.Collection(file.Collection).Drop(ctx); err != nil {
			log.Printf("Błąd przy czyszczeniu kolekcji %s: %v", file.Collection, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd czyszczenia kolekcji " + file.Collection})
			return
		}

		docs := toInterfaceSlice(file.Model)
		if len(docs) > 0 {
			_, err := db.Collection(file.Collection).InsertMany(ctx, docs)
			if err != nil {
				log.Printf("Błąd przy wstawianiu do %s: %v", file.Collection, err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd wstawiania do kolekcji " + file.Collection})
				return
			}
		}
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
	case *[]models.Order:
		for _, item := range *v {
			res = append(res, item)
		}
	case *[]models.Cart:
		for _, item := range *v {
			res = append(res, item)
		}
	case *[]models.Review:
		for _, item := range *v {
			res = append(res, item)
		}
	}
	return res
}

func setTimestamps(data interface{}) {
	now := time.Now()
	switch v := data.(type) {
	case *[]models.Album:
		for i := range *v {
			(*v)[i].CreatedAt = now
			(*v)[i].UpdatedAt = now
		}
	case *[]models.User:
		for i := range *v {
			(*v)[i].CreatedAt = now
			(*v)[i].UpdatedAt = now
		}
	case *[]models.Order:
		for i := range *v {
			(*v)[i].CreatedAt = now
			(*v)[i].UpdatedAt = now
		}
	case *[]models.Review:
		for i := range *v {
			(*v)[i].CreatedAt = now
		}
	}
}
