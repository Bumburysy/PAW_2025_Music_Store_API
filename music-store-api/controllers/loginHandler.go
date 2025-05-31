package controllers

import (
	"context"
	"log"
	"music-store-api/config"
	"music-store-api/middleware"
	"music-store-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// LoginRequest reprezentuje payload do logowania
type LoginRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"strongpassword"`
}

// LoginResponse reprezentuje odpowiedź po zalogowaniu (token JWT)
type LoginResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// Login godoc
// @Summary Logowanie użytkownika
// @Description Zwraca token JWT po poprawnym zalogowaniu
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "Dane logowania"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} models.ErrorResponse "Niepoprawne dane"
// @Failure 401 {object} models.ErrorResponse "Błędne dane logowania"
// @Router /login [post]
func Login(c *gin.Context) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Niepoprawne dane"})
		return
	}

	var user models.User
	err := config.DB.Collection("users").FindOne(context.Background(), bson.M{"email": credentials.Email}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Nieprawidłowy email lub hasło"})
		return
	}

	if !middleware.CheckPasswordHash(credentials.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Nieprawidłowy email lub hasło"})
		log.Printf("User found: %s, PasswordHash: %s", user.Email, user.PasswordHash)
		return
	}

	token, err := config.GenerateJWT(user.ID.Hex(), user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Błąd generowania tokena"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
