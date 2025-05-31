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
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
	}
	return res
}
