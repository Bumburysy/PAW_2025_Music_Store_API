package controllers

import (
	"context"
	"music-store-api/config"
	"music-store-api/middleware"
	"music-store-api/models"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection

// InitUserCollection inicjalizuje kolekcję użytkowników
func InitUserCollection() {
	userCollection = config.DB.Collection("users")
}

// GetUsers godoc
// @Summary Pobierz listę użytkowników
// @Security BearerAuth
// @Description Zwraca wszystkich użytkowników w systemie
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {array} models.User
// @Router /users [get]
func GetUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := userCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd pobierania użytkowników"})
		return
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err = cursor.All(ctx, &users); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd dekodowania użytkowników"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// GetUserByID godoc
// @Summary Pobierz użytkownika po ID
// @Security BearerAuth
// @Description Zwraca szczegóły użytkownika na podstawie ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "ID użytkownika"
// @Success 200 {object} models.User
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /users/{id} [get]
func GetUserByID(c *gin.Context) {
	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Niepoprawne ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err = userCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Użytkownik nie znaleziony"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// CreateUser godoc
// @Summary Dodaj nowego użytkownika
// @Security BearerAuth
// @Description Tworzy nowego użytkownika w bazie
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.User true "Użytkownik do dodania"
// @Success 201 {object} models.User
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /users [post]
func CreateUser(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Niepoprawne dane wejściowe"})
		return
	}

	hashedPassword, err := middleware.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd haszowania hasła"})
		return
	}
	user.PasswordHash = hashedPassword

	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = userCollection.InsertOne(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd tworzenia użytkownika"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"_id": user.ID.Hex()})
}

// UpdateUser godoc
// @Summary Aktualizuj użytkownika
// @Security BearerAuth
// @Description Aktualizuje dane użytkownika na podstawie ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "ID użytkownika"
// @Param user body models.User true "Zaktualizowane dane użytkownika"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /users/{id} [patch]
func UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Niepoprawne ID"})
		return
	}

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Niepoprawne dane wejściowe"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"email":      user.Email,
			"role":       user.Role,
			"password":   user.Password,
			"updated_at": time.Now(),
		},
	}

	if user.Password != "" {
		hashedPassword, err := middleware.HashPassword(user.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd haszowania hasła"})
			return
		}
		update["$set"].(bson.M)["password_hash"] = hashedPassword
	}

	result, err := userCollection.UpdateByID(ctx, objID, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd aktualizacji użytkownika"})
		return
	}
	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Użytkownik nie znaleziony"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Użytkownik zaktualizowany"})
}

// DeleteUser godoc
// @Summary Usuń użytkownika
// @Security BearerAuth
// @Description Usuwa użytkownika na podstawie ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "ID użytkownika"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /users/{id} [delete]
func DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Niepoprawne ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := userCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd usuwania użytkownika"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Użytkownik nie znaleziony"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Użytkownik usunięty"})
}
