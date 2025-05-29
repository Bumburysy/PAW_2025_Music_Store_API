package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"music-store-api/config"
	"music-store-api/models"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func LoadTestData(c *gin.Context) {
	ctx := context.Background()
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
		path := fmt.Sprintf("data/%s", file.Filename)
		byteValue, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("Nie udało się otworzyć %s: %v\n", path, err)
			continue
		}

		err = json.Unmarshal(byteValue, file.Model)
		if err != nil {
			fmt.Printf("Błąd dekodowania %s: %v\n", path, err)
			continue
		}

		if err := db.Collection(file.Collection).Drop(ctx); err != nil {
			fmt.Printf("Błąd przy czyszczeniu kolekcji %s: %v\n", file.Collection, err)
			continue
		}

		docs := toInterfaceSlice(file.Model)
		if len(docs) > 0 {
			_, err := db.Collection(file.Collection).InsertMany(ctx, docs)
			if err != nil {
				fmt.Printf("Błąd przy wstawianiu do %s: %v\n", file.Collection, err)
				continue
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
