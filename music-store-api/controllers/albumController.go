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

// GetAlbums godoc
// @Summary Pobierz listę albumów
// @Description Zwraca wszystkie albumy w sklepie
// @Tags albums
// @Accept json
// @Produce json
// @Success 200 {array} models.Album
// @Router /albums [get]
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

// GetAlbumByID godoc
// @Summary Pobierz album po ID
// @Description Zwraca szczegóły albumu na podstawie ID
// @Tags albums
// @Accept json
// @Produce json
// @Param id path string true "ID albumu"
// @Success 200 {object} models.Album
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /albums/{id} [get]
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

// CreateAlbum godoc
// @Summary Dodaj nowy album
// @Description Dodaje album do bazy danych
// @Tags albums
// @Accept json
// @Produce json
// @Param album body models.Album true "Album do dodania"
// @Success 201 {object} models.Album
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /albums [post]
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

// UpdateAlbum godoc
// @Summary Zaktualizuj album
// @Description Aktualizuje dane albumu na podstawie ID
// @Tags albums
// @Accept json
// @Produce json
// @Param id path string true "ID albumu"
// @Param album body models.Album true "Zaktualizowane dane albumu"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /albums/{id} [patch]
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

// DeleteAlbum godoc
// @Summary Usuń album
// @Description Usuwa album na podstawie ID
// @Tags albums
// @Accept json
// @Produce json
// @Param id path string true "ID albumu"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /albums/{id} [delete]
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
