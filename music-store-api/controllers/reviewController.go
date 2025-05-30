package controllers

import (
	"context"
	"music-store-api/config"
	"music-store-api/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var reviewCollection *mongo.Collection

func InitReviewCollection() {
	reviewCollection = config.DB.Collection("reviews")
}

// GetReviews godoc
// @Summary Pobierz wszystkie recenzje
// @Tags Reviews
// @Produce json
// @Success 200 {array} models.Review
// @Failure 500 {object} map[string]string
// @Router /reviews [get]
func GetReviews(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := reviewCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd pobierania recenzji"})
		return
	}
	defer cursor.Close(ctx)

	var reviews []models.Review
	if err = cursor.All(ctx, &reviews); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd dekodowania danych"})
		return
	}

	c.JSON(http.StatusOK, reviews)
}

// GetReviewByID godoc
// @Summary Pobierz recenzję po ID
// @Tags Reviews
// @Produce json
// @Param id path string true "ID recenzji"
// @Success 200 {object} models.Review
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /reviews/{id} [get]
func GetReviewByID(c *gin.Context) {
	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Niepoprawne ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var review models.Review
	err = reviewCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&review)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recenzja nie znaleziona"})
		return
	}

	c.JSON(http.StatusOK, review)
}

// GetReviewsByAlbumID godoc
// @Summary Pobierz recenzje dla konkretnego albumu
// @Tags Reviews
// @Produce json
// @Param albumID path string true "ID albumu"
// @Success 200 {array} models.Review
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /reviews/album/{albumID} [get]
func GetReviewsByAlbumID(c *gin.Context) {
	albumIDParam := c.Param("albumID")
	albumID, err := primitive.ObjectIDFromHex(albumIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Niepoprawne ID albumu"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := reviewCollection.Find(ctx, bson.M{"album_id": albumID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd pobierania recenzji"})
		return
	}
	defer cursor.Close(ctx)

	var reviews []models.Review
	if err = cursor.All(ctx, &reviews); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd dekodowania danych"})
		return
	}

	c.JSON(http.StatusOK, reviews)
}

// GetReviewsByUserID godoc
// @Summary Pobierz recenzje użytkownika
// @Tags Reviews
// @Produce json
// @Param userID path string true "ID użytkownika"
// @Success 200 {array} models.Review
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /reviews/user/{userID} [get]
func GetReviewsByUserID(c *gin.Context) {
	userIDParam := c.Param("userID")
	userID, err := primitive.ObjectIDFromHex(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Niepoprawne ID użytkownika"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := reviewCollection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd pobierania recenzji"})
		return
	}
	defer cursor.Close(ctx)

	var reviews []models.Review
	if err = cursor.All(ctx, &reviews); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd dekodowania danych"})
		return
	}

	c.JSON(http.StatusOK, reviews)
}

// CreateReview godoc
// @Summary Dodaj nową recenzję
// @Tags Reviews
// @Accept json
// @Produce json
// @Param review body models.Review true "Nowa recenzja"
// @Success 201 {object} models.Review
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /reviews [post]
func CreateReview(c *gin.Context) {
	var review models.Review
	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Niepoprawne dane"})
		return
	}

	if review.Rating < 1 || review.Rating > 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ocena musi być w zakresie 1-5"})
		return
	}

	if review.AlbumID.IsZero() || review.UserID.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "AlbumID i UserID są wymagane"})
		return
	}

	review.ID = primitive.NewObjectID()
	review.CreatedAt = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := reviewCollection.InsertOne(ctx, review)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd dodawania recenzji"})
		return
	}

	c.JSON(http.StatusCreated, review)
}

// UpdateReview godoc
// @Summary Zaktualizuj recenzję
// @Tags Reviews
// @Accept json
// @Produce json
// @Param id path string true "ID recenzji"
// @Param review body models.Review true "Dane recenzji do aktualizacji"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /reviews/{id} [put]
func UpdateReview(c *gin.Context) {
	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Niepoprawne ID"})
		return
	}

	var review models.Review
	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Niepoprawne dane"})
		return
	}

	if review.Rating < 1 || review.Rating > 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ocena musi być w zakresie 1-5"})
		return
	}

	if review.AlbumID.IsZero() || review.UserID.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "AlbumID i UserID są wymagane"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"album_id": review.AlbumID,
			"user_id":  review.UserID,
			"rating":   review.Rating,
			"comment":  review.Comment,
		},
	}

	result, err := reviewCollection.UpdateByID(ctx, objID, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd aktualizacji recenzji"})
		return
	}
	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recenzja nie znaleziona"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Recenzja zaktualizowana"})
}

// DeleteReview godoc
// @Summary Usuń recenzję
// @Tags Reviews
// @Produce json
// @Param id path string true "ID recenzji"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /reviews/{id} [delete]
func DeleteReview(c *gin.Context) {
	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Niepoprawne ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := reviewCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd usuwania recenzji"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recenzja nie znaleziona"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Recenzja usunięta"})
}
