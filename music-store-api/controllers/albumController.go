package controllers

import (
	"context"
	"net/http"
	"time"

	"music-store-api/config"
	"music-store-api/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Kolekcja albumów z MongoDB
var albumCollection *mongo.Collection

func InitAlbumCollection() {
	albumCollection = config.DB.Collection("albums")
}

func GetAlbums(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := albumCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd pobierania albumów"})
		return
	}
	defer cursor.Close(ctx)

	var albums []models.Album
	if err = cursor.All(ctx, &albums); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd dekodowania albumów"})
		return
	}

	c.JSON(http.StatusOK, albums)
}

func GetAlbumByID(c *gin.Context) {
	idParam := c.Param("id")

	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Niepoprawne ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var album models.Album
	err = albumCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&album)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Album nie znaleziony"})
		return
	}

	c.JSON(http.StatusOK, album)
}

func CreateAlbum(c *gin.Context) {
	var album models.Album

	if err := c.ShouldBindJSON(&album); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Niepoprawne dane wejściowe"})
		return
	}

	album.ID = primitive.NewObjectID()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := albumCollection.InsertOne(ctx, album)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd tworzenia albumu"})
		return
	}

	c.JSON(http.StatusCreated, album)
}

func UpdateAlbum(c *gin.Context) {
	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Niepoprawne ID"})
		return
	}

	var album models.Album
	if err := c.ShouldBindJSON(&album); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Niepoprawne dane wejściowe"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"title":  album.Title,
			"artist": album.Artist,
			"price":  album.Price,
			"genre":  album.Genre,
			"stock":  album.Stock,
		},
	}

	result, err := albumCollection.UpdateByID(ctx, objID, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd aktualizacji albumu"})
		return
	}
	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Album nie znaleziony"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Album zaktualizowany"})
}

func DeleteAlbum(c *gin.Context) {
	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Niepoprawne ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := albumCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd usuwania albumu"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Album nie znaleziony"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Album usunięty"})
}
