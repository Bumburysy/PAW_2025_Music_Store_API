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

var cartCollection *mongo.Collection

func InitCartCollection() {
	cartCollection = config.DB.Collection("carts")
}

// GetCarts godoc
// @Summary Pobierz wszystkie koszyki
// @Tags Carts
// @Produce json
// @Success 200 {array} models.Cart
// @Failure 500 {object} map[string]string
// @Router /carts [get]
func GetCarts(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := cartCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd pobierania koszyków"})
		return
	}
	defer cursor.Close(ctx)

	var carts []models.Cart
	if err = cursor.All(ctx, &carts); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd dekodowania danych"})
		return
	}

	c.JSON(http.StatusOK, carts)
}

// GetCartByID godoc
// @Summary Pobierz koszyk po ID
// @Tags Carts
// @Produce json
// @Param id path string true "ID koszyka"
// @Success 200 {object} models.Cart
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /carts/{id} [get]
func GetCartByID(c *gin.Context) {
	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Niepoprawne ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var cart models.Cart
	err = cartCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&cart)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Koszyk nie znaleziony"})
		return
	}

	c.JSON(http.StatusOK, cart)
}

// GetCartByUserID godoc
// @Summary Pobierz koszyk użytkownika
// @Tags Carts
// @Produce json
// @Param userID path string true "ID użytkownika"
// @Success 200 {object} models.Cart
// @Failure 404 {object} map[string]string
// @Router /carts/user/{userID} [get]
func GetCartByUserID(c *gin.Context) {
	userID := c.Param("userID")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var cart models.Cart
	err := cartCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&cart)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Koszyk nie znaleziony"})
		return
	}

	c.JSON(http.StatusOK, cart)
}

// CreateCart godoc
// @Summary Utwórz nowy koszyk
// @Tags Carts
// @Accept json
// @Produce json
// @Param cart body models.Cart true "Nowy koszyk"
// @Success 201 {object} models.Cart
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /carts [post]
func CreateCart(c *gin.Context) {
	var cart models.Cart
	if err := c.ShouldBindJSON(&cart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Niepoprawne dane"})
		return
	}

	cart.ID = primitive.NewObjectID()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := cartCollection.InsertOne(ctx, cart)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd tworzenia koszyka"})
		return
	}

	c.JSON(http.StatusCreated, cart)
}

// UpdateCart godoc
// @Summary Zaktualizuj koszyk
// @Tags Carts
// @Accept json
// @Produce json
// @Param id path string true "ID koszyka"
// @Param cart body models.Cart true "Dane koszyka"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /carts/{id} [put]
func UpdateCart(c *gin.Context) {
	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Niepoprawne ID"})
		return
	}

	var cart models.Cart
	if err := c.ShouldBindJSON(&cart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Niepoprawne dane"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{
		"$set": cart,
	}

	result, err := cartCollection.UpdateByID(ctx, objID, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd aktualizacji"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Koszyk nie znaleziony"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Koszyk zaktualizowany"})
}

// DeleteCart godoc
// @Summary Usuń koszyk
// @Tags Carts
// @Produce json
// @Param id path string true "ID koszyka"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /carts/{id} [delete]
func DeleteCart(c *gin.Context) {
	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Niepoprawne ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := cartCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd usuwania koszyka"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Koszyk nie znaleziony"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Koszyk usunięty"})
}

// AddItemToCart godoc
// @Summary Dodaj pozycję do koszyka lub zaktualizuj ilość
// @Tags Carts
// @Accept json
// @Produce json
// @Param cartID path string true "ID koszyka"
// @Param item body models.CartItem true "Pozycja do dodania"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /carts/{cartID}/items [post]
func AddItemToCart(c *gin.Context) {
	cartID := c.Param("cartID")
	objID, err := primitive.ObjectIDFromHex(cartID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Niepoprawne ID koszyka"})
		return
	}

	var item models.CartItem
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Niepoprawne dane pozycji"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": objID, "items.album_id": item.AlbumID}
	update := bson.M{"$inc": bson.M{"items.$.quantity": item.Quantity}}
	result, err := cartCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd dodawania pozycji"})
		return
	}

	if result.MatchedCount == 0 {
		update = bson.M{"$push": bson.M{"items": item}}
		_, err = cartCollection.UpdateByID(ctx, objID, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd dodawania nowej pozycji"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pozycja dodana lub zaktualizowana", "item": item})
}

// RemoveItemFromCart godoc
// @Summary Usuń pozycję z koszyka
// @Tags Carts
// @Produce json
// @Param cartID path string true "ID koszyka"
// @Param albumID path string true "ID albumu"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /carts/{cartID}/items/{albumID} [delete]
func RemoveItemFromCart(c *gin.Context) {
	cartID := c.Param("cartID")
	albumIDParam := c.Param("albumID")
	albumID, err := primitive.ObjectIDFromHex(albumIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Niepoprawne ID albumu"})
		return
	}
	objID, err := primitive.ObjectIDFromHex(cartID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Niepoprawne ID koszyka"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{
		"$pull": bson.M{"items": bson.M{"album_id": albumID}},
	}

	_, err = cartCollection.UpdateByID(ctx, objID, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd usuwania pozycji"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pozycja usunięta"})
}

// UpdateCartItemQuantity godoc
// @Summary Zaktualizuj ilość pozycji w koszyku
// @Tags Carts
// @Accept json
// @Produce json
// @Param cartID path string true "ID koszyka"
// @Param albumID path string true "ID albumu"
// @Param quantity body map[string]string true "Nowa ilość"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /carts/{cartID}/items/{albumID} [put]
func UpdateCartItemQuantity(c *gin.Context) {
	cartID := c.Param("cartID")
	albumIDParam := c.Param("albumID")
	albumID, err := primitive.ObjectIDFromHex(albumIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Niepoprawne ID albumu"})
		return
	}
	objID, err := primitive.ObjectIDFromHex(cartID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Niepoprawne ID koszyka"})
		return
	}

	var payload struct {
		Quantity int `json:"quantity"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Niepoprawne dane"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": objID, "items.album_id": albumID}
	update := bson.M{"$set": bson.M{"items.$.quantity": payload.Quantity}}

	result, err := cartCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd aktualizacji"})
		return
	}
	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pozycja nie znaleziona"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ilość zaktualizowana"})
}

// UpdateCartTotal godoc
// @Summary Oblicz i zaktualizuj sumę wartości koszyka
// @Tags Carts
// @Produce json
// @Param cartID path string true "ID koszyka"
// @Success 200 {object} models.Cart
// @Failure 400 {object} map[string]string "Niepoprawne ID koszyka"
// @Failure 404 {object} map[string]string "Koszyk nie znaleziony"
// @Failure 500 {object} map[string]string "Błąd serwera podczas pobierania lub aktualizacji danych"
// @Router /carts/{cartID}/total [put]
func UpdateCartTotal(c *gin.Context) {
	cartID := c.Param("cartID")
	oid, err := primitive.ObjectIDFromHex(cartID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Niepoprawne ID koszyka"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var cart models.Cart
	err = cartCollection.FindOne(ctx, bson.M{"_id": oid}).Decode(&cart)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Koszyk nie znaleziony"})
		return
	}

	total := 0.0
	for _, item := range cart.Items {
		var album models.Album
		err = albumCollection.FindOne(ctx, bson.M{"_id": item.AlbumID}).Decode(&album)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd pobierania danych albumu"})
			return
		}
		total += album.Price * float64(item.Quantity)
	}

	update := bson.M{"$set": bson.M{"total": total}}
	_, err = cartCollection.UpdateOne(ctx, bson.M{"_id": oid}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd aktualizacji sumy koszyka"})
		return
	}

	cart.Total = total
	c.JSON(http.StatusOK, cart)
}

// ClearCart godoc
// @Summary Wyczyść koszyk (usuń wszystkie pozycje i zresetuj sumę)
// @Tags Carts
// @Produce json
// @Param id path string true "ID koszyka"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /carts/{id}/clear [post]
func ClearCart(c *gin.Context) {
	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Niepoprawne ID koszyka"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"items": []models.CartItem{},
			"total": 0,
		},
	}

	result, err := cartCollection.UpdateByID(ctx, objID, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd czyszczenia koszyka"})
		return
	}
	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Koszyk nie znaleziony"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Koszyk wyczyszczony"})
}
