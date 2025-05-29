package controllers

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"music-store-api/config"
	"music-store-api/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Kolekcja albumów z MongoDB
var albumCollection *mongo.Collection

func InitAlbumCollection() {
	albumCollection = config.DB.Collection("albums")
}

// GetAlbums godoc
// @Summary Pobierz listę albumów
// @Description Zwraca wszystkie albumy w sklepie z opcjonalnym filtrowaniem, sortowaniem i paginacją
// @Tags Albums
// @Accept json
// @Produce json
// @Param page query int false "Numer strony (domyślnie 1)"
// @Param limit query int false "Liczba wyników na stronę (domyślnie 10)"
// @Param artist query string false "Filtruj po wykonawcy (częściowa zgodność, bez wielkości liter)"
// @Param genre query string false "Filtruj po gatunku muzycznym (częściowa zgodność, bez wielkości liter)"
// @Param sort query string false "Sortowanie po polach (np. price,-title)"
// @Success 200 {object} map[string]interface{} "Struktura danych zawiera: page, limit, total i data (lista albumów)"
// @Failure 500 {object} map[string]string
// @Router /albums [get]
func GetAlbums(c *gin.Context) {
	// Parsujemy parametry zapytania
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}
	artist := c.Query("artist")
	genre := c.Query("genre")
	sort := c.DefaultQuery("sort", "")

	// Budujemy filtr
	filter := bson.M{}
	if artist != "" {
		filter["artist"] = bson.M{"$regex": artist, "$options": "i"}
	}
	if genre != "" {
		filter["genre"] = bson.M{"$regex": genre, "$options": "i"}
	}

	// Ustawienia sortowania
	findOptions := options.Find()
	if sort != "" {
		sortFields := bson.D{}
		for _, field := range strings.Split(sort, ",") {
			direction := 1
			if strings.HasPrefix(field, "-") {
				direction = -1
				field = field[1:]
			}
			sortFields = append(sortFields, bson.E{Key: field, Value: direction})
		}
		findOptions.SetSort(sortFields)
	}

	// Paginacja (skip i limit)
	findOptions.SetSkip(int64((page - 1) * limit))
	findOptions.SetLimit(int64(limit))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Pobieramy dane
	cursor, err := albumCollection.Find(ctx, filter, findOptions)
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

	// Liczymy liczbę wszystkich dokumentów dla podanego filtra (do podania total count)
	total, _ := albumCollection.CountDocuments(ctx, filter)

	c.JSON(http.StatusOK, gin.H{
		"page":  page,
		"limit": limit,
		"total": total,
		"data":  albums,
	})
}

// GetAlbumByID godoc
// @Summary Pobierz album po ID
// @Description Zwraca szczegóły albumu na podstawie ID
// @Tags Albums
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
// @Tags Albums
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

// CreateAlbumsBulk godoc
// @Summary Dodaj wiele albumów naraz
// @Description Dodaje wiele albumów do bazy danych w jednym żądaniu
// @Tags Albums
// @Accept json
// @Produce json
// @Param albums body []models.Album true "Lista albumów do dodania"
// @Success 201 {object} map[string]interface{} "Informacja o dodanych albumach i ich liczbie"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /albums/bulk [post]
func CreateAlbumsBulk(c *gin.Context) {
	var albums []models.Album

	if err := c.ShouldBindJSON(&albums); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Niepoprawne dane wejściowe"})
		return
	}

	var docs []interface{}
	for _, album := range albums {
		album.ID = primitive.NewObjectID()
		docs = append(docs, album)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := albumCollection.InsertMany(ctx, docs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd przy dodawaniu albumów"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Albumy zostały dodane", "count": len(albums)})
}

// UpdateAlbum godoc
// @Summary Zaktualizuj album
// @Description Aktualizuje dane albumu na podstawie ID
// @Tags Albums
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
			"title":    album.Title,
			"artist":   album.Artist,
			"price":    album.Price,
			"genre":    album.Genre,
			"quantity": album.Quantity,
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
// @Tags Albums
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
