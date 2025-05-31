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

var orderCollection *mongo.Collection

func InitOrderCollection() {
	orderCollection = config.DB.Collection("orders")
}

// GetOrders godoc
// @Summary Pobierz wszystkie zamówienia
// @Security BearerAuth
// @Tags Orders
// @Produce json
// @Success 200 {array} models.Order
// @Failure 500 {object} models.ErrorResponse
// @Router /orders [get]
func GetOrders(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := orderCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd pobierania zamówień"})
		return
	}
	defer cursor.Close(ctx)

	var orders []models.Order
	if err = cursor.All(ctx, &orders); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd dekodowania danych"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// GetOrderByID godoc
// @Summary Pobierz zamówienie po ID
// @Security BearerAuth
// @Tags Orders
// @Produce json
// @Param id path string true "ID zamówienia"
// @Success 200 {object} models.Order
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /orders/{id} [get]
func GetOrderByID(c *gin.Context) {
	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Niepoprawne ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var order models.Order
	err = orderCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&order)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Zamówienie nie znalezione"})
		return
	}

	c.JSON(http.StatusOK, order)
}

// GetOrdersByUserID godoc
// @Summary Pobierz zamówienia użytkownika
// @Security BearerAuth
// @Tags Orders
// @Produce json
// @Param userID path string true "ID użytkownika"
// @Success 200 {array} models.Order
// @Failure 404 {object} models.ErrorResponse
// @Router /orders/user/{userID} [get]
func GetOrdersByUserID(c *gin.Context) {
	userIDParam := c.Param("userID")
	userID, err := primitive.ObjectIDFromHex(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Niepoprawne ID użytkownika"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := orderCollection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd pobierania zamówień"})
		return
	}
	defer cursor.Close(ctx)

	var orders []models.Order
	if err = cursor.All(ctx, &orders); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd dekodowania danych"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// CreateOrder godoc
// @Summary Utwórz nowe zamówienie
// @Security BearerAuth
// @Tags Orders
// @Accept json
// @Produce json
// @Param order body models.Order true "Nowe zamówienie"
// @Success 201 {object} models.Order
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /orders [post]
func CreateOrder(c *gin.Context) {
	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Niepoprawne dane"})
		return
	}

	order.ID = primitive.NewObjectID()
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := orderCollection.InsertOne(ctx, order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd tworzenia zamówienia"})
		return
	}

	c.JSON(http.StatusCreated, order)
}

// UpdateOrder godoc
// @Summary Zaktualizuj zamówienie (np. status)
// @Security BearerAuth
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "ID zamówienia"
// @Param order body models.Order true "Dane zamówienia do aktualizacji"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /orders/{id} [put]
func UpdateOrder(c *gin.Context) {
	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Niepoprawne ID"})
		return
	}

	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Niepoprawne dane"})
		return
	}

	order.UpdatedAt = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{
		"$set": order,
	}

	result, err := orderCollection.UpdateByID(ctx, objID, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd aktualizacji"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Zamówienie nie znalezione"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Zamówienie zaktualizowane"})
}

// DeleteOrder godoc
// @Summary Usuń zamówienie
// @Security BearerAuth
// @Tags Orders
// @Produce json
// @Param id path string true "ID zamówienia"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /orders/{id} [delete]
func DeleteOrder(c *gin.Context) {
	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Niepoprawne ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := orderCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd usuwania zamówienia"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Zamówienie nie znalezione"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Zamówienie usunięte"})
}

// UpdateOrderStatus godoc
// @Summary Zaktualizuj status zamówienia
// @Security BearerAuth
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "ID zamówienia"
// @Param status body map[string]string true "Nowy status zamówienia"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /orders/{id}/status [patch]
func UpdateOrderStatus(c *gin.Context) {
	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Niepoprawne ID"})
		return
	}

	var body struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.Status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Niepoprawny status"})
		return
	}

	if body.Status != models.OrderStatusPending &&
		body.Status != models.OrderStatusProcessing &&
		body.Status != models.OrderStatusShipped &&
		body.Status != models.OrderStatusCompleted &&
		body.Status != models.OrderStatusCancelled {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Niepoprawny status"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"status":     body.Status,
			"updated_at": time.Now(),
		},
	}

	result, err := orderCollection.UpdateByID(ctx, objID, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd aktualizacji statusu"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Zamówienie nie znalezione"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status zamówienia zaktualizowany"})
}

// UpdateOrderShipping godoc
// @Summary Zaktualizuj dane wysyłki zamówienia
// @Security BearerAuth
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "ID zamówienia"
// @Param shipping body models.ShippingDetails true "Nowe dane wysyłki"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /orders/{id}/shipping [put]
func UpdateOrderShipping(c *gin.Context) {
	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Niepoprawne ID"})
		return
	}

	var shipping models.ShippingDetails
	if err := c.ShouldBindJSON(&shipping); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Niepoprawne dane wysyłki"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"shipping":   shipping,
			"updated_at": time.Now(),
		},
	}

	result, err := orderCollection.UpdateByID(ctx, objID, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd aktualizacji danych wysyłki"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Zamówienie nie znalezione"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Dane wysyłki zaktualizowane"})
}
