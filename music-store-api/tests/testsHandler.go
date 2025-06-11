package tests

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"

	"music-store-api/controllers"
	"music-store-api/middleware"

	"github.com/gin-gonic/gin"
)

func RunTestsHandler(c *gin.Context) {
	log.Println("Rozpoczynam testy API...")

	// Login i pobranie tokena
	token, err := loginAndGetToken()
	if err != nil {
		log.Printf("Błąd logowania: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Login failed: " + err.Error()})
		return
	}
	log.Println("Pomyślnie zalogowano i pobrano token.")

	// Testy ogólne i albumów
	err = RunBasicApiTests(token)
	if err != nil {
		log.Printf("Błąd testów API: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Println("Wszystkie testy zakończone sukcesem.")
	c.JSON(http.StatusOK, gin.H{"result": "success"})
}

// Login admina i pobranie tokena
func loginAndGetToken() (string, error) {
	router := SetupTestRouter()
	loginPayload := `{"email":"jan.kowalski@example.com","password":"password1"}`

	req, _ := http.NewRequest("POST", "/login", strings.NewReader(loginPayload))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		return "", fmt.Errorf("POST /login failed with status %d", resp.Code)
	}

	var loginResp struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(resp.Body.Bytes(), &loginResp); err != nil {
		return "", fmt.Errorf("token parsing error: %v", err)
	}
	if loginResp.Token == "" {
		return "", fmt.Errorf("empty token returned")
	}

	return loginResp.Token, nil
}

// Wszystkie testy z tokenem
func RunBasicApiTests(token string) error {
	router := SetupTestRouter()

	// CRUD /albums

	// 1. GET /albums (publiczny)
	log.Println("Test: GET /albums (publiczny)")
	req1, _ := http.NewRequest("GET", "/albums", nil)
	resp1 := httptest.NewRecorder()
	router.ServeHTTP(resp1, req1)
	if resp1.Code != http.StatusOK {
		return fmt.Errorf("GET /albums failed with status %d", resp1.Code)
	}
	log.Println("OK: GET /albums")

	// 2. POST /albums bez tokena (ma się nie udać)
	log.Println("Test: POST /albums bez tokena (powinno zwrócić 401)")
	albumJSON := `{
    "title": "The Dark Side of the Moon",
    "artist": "Pink Floyd",
    "genre": "Progressive Rock",
    "description": "Classic album from Pink Floyd.",
    "release_date": "1973-03-01T00:00:00Z",
    "tracks": ["Speak to Me", "Breathe", "Time", "Money"],
    "price": 29.99,
    "quantity": 100,
    "cover_url": "https://example.com/darkside.jpg"}`
	req2, _ := http.NewRequest("POST", "/albums", strings.NewReader(albumJSON))
	req2.Header.Set("Content-Type", "application/json")
	resp2 := httptest.NewRecorder()
	router.ServeHTTP(resp2, req2)
	if resp2.Code != http.StatusUnauthorized {
		return fmt.Errorf("POST /albums without auth expected 401, got %d", resp2.Code)
	}
	log.Println("OK: POST /albums bez tokena (401)")

	// 3. POST /albums z tokenem
	log.Println("Test: POST /albums z tokenem")
	req3, _ := http.NewRequest("POST", "/albums", strings.NewReader(albumJSON))
	req3.Header.Set("Content-Type", "application/json")
	req3.Header.Set("Authorization", "Bearer "+token)
	resp3 := httptest.NewRecorder()
	router.ServeHTTP(resp3, req3)
	if resp3.Code != http.StatusCreated {
		return fmt.Errorf("POST /albums with auth expected 201, got %d", resp3.Code)
	}
	log.Println("OK: POST /albums z tokenem (201)")

	// 4. Parsuj ID nowego albumu
	var createdAlbum struct {
		ID string `json:"_id"`
	}
	if err := json.Unmarshal(resp3.Body.Bytes(), &createdAlbum); err != nil {
		return fmt.Errorf("parsing created album failed: %v", err)
	}
	if createdAlbum.ID == "" {
		return fmt.Errorf("album ID is empty")
	}
	log.Printf("Utworzono album o ID: %s\n", createdAlbum.ID)

	// 5. GET /albums/:id
	log.Println("Test: GET /albums/:id")
	req4, _ := http.NewRequest("GET", "/albums/"+createdAlbum.ID, nil)
	resp4 := httptest.NewRecorder()
	router.ServeHTTP(resp4, req4)
	if resp4.Code != http.StatusOK {
		return fmt.Errorf("GET /albums/:id failed: expected 200, got %d", resp4.Code)
	}
	log.Println("OK: GET /albums/:id")

	// 6. PATCH /albums/:id
	log.Println("Test: PATCH /albums/:id")
	updateAlbumJSON := `{"title": "Zmieniony Tytuł"}`
	req5, _ := http.NewRequest("PATCH", "/albums/"+createdAlbum.ID, strings.NewReader(updateAlbumJSON))
	req5.Header.Set("Content-Type", "application/json")
	req5.Header.Set("Authorization", "Bearer "+token)
	resp5 := httptest.NewRecorder()
	router.ServeHTTP(resp5, req5)
	if resp5.Code != http.StatusOK {
		return fmt.Errorf("PATCH /albums/:id failed: expected 200, got %d", resp5.Code)
	}
	log.Println("OK: PATCH /albums/:id")

	// 7. DELETE /albums/:id
	log.Println("Test: DELETE /albums/:id")
	req6, _ := http.NewRequest("DELETE", "/albums/"+createdAlbum.ID, nil)
	req6.Header.Set("Authorization", "Bearer "+token)
	resp6 := httptest.NewRecorder()
	router.ServeHTTP(resp6, req6)
	if resp6.Code != http.StatusOK {
		return fmt.Errorf("DELETE /albums/:id failed: expected 200, got %d", resp6.Code)
	}
	log.Println("OK: DELETE /albums/:id")

	// CRUD /users

	// 1. POST /users (utwórz nowego użytkownika)
	log.Println("Test: POST /users (tworzenie użytkownika)")
	userJSON := `{
    "first_name": "Adam",
    "last_name": "Kowalski",
    "email": "adam.kowalski@example.com",
    "phone_number": "+48123434589",
    "password": "password12",
    "role": "admin",
    "is_active": true,
    "shipping_details": {
      "address": "ul. Kwiatowa 15",
      "city": "Warszawa",
      "postal_code": "00-001",
      "country": "Polska",
      "phone_number": "+48123326789"
        }
    }`
	req7, _ := http.NewRequest("POST", "/users", strings.NewReader(userJSON))
	req7.Header.Set("Content-Type", "application/json")
	req7.Header.Set("Authorization", "Bearer "+token)
	resp7 := httptest.NewRecorder()
	router.ServeHTTP(resp7, req7)
	if resp7.Code != http.StatusCreated {
		return fmt.Errorf("POST /users failed: expected 201, got %d", resp7.Code)
	}
	log.Println("OK: POST /users (201)")

	// 2. Parsuj ID użytkownika
	var createdUser struct {
		ID string `json:"_id"`
	}
	if err := json.Unmarshal(resp7.Body.Bytes(), &createdUser); err != nil {
		return fmt.Errorf("parsing created user failed: %v", err)
	}
	if createdUser.ID == "" {
		return fmt.Errorf("user ID is empty after creation")
	}
	log.Printf("Utworzono użytkownika o ID: %s\n", createdUser.ID)

	// 3. GET /users (lista)
	log.Println("Test: GET /users (lista)")
	req8, _ := http.NewRequest("GET", "/users", nil)
	req8.Header.Set("Authorization", "Bearer "+token)
	resp8 := httptest.NewRecorder()
	router.ServeHTTP(resp8, req8)
	if resp8.Code != http.StatusOK {
		return fmt.Errorf("GET /users failed: expected 200, got %d", resp8.Code)
	}
	log.Println("OK: GET /users")

	// 4. GET /users/:id
	log.Println("Test: GET /users/:id")
	req9, _ := http.NewRequest("GET", "/users/"+createdUser.ID, nil)
	req9.Header.Set("Authorization", "Bearer "+token)
	resp9 := httptest.NewRecorder()
	router.ServeHTTP(resp9, req9)
	if resp9.Code != http.StatusOK {
		return fmt.Errorf("GET /users/:id failed: expected 200, got %d", resp9.Code)
	}
	log.Println("OK: GET /users/:id")

	// 5. PATCH /users/:id (aktualizacja imienia)
	log.Println("Test: PATCH /users/:id (zmiana imienia)")
	updateUserJSON := `{"first_name": "Zmieniona"}`
	req10, _ := http.NewRequest("PATCH", "/users/"+createdUser.ID, strings.NewReader(updateUserJSON))
	req10.Header.Set("Content-Type", "application/json")
	req10.Header.Set("Authorization", "Bearer "+token)
	resp10 := httptest.NewRecorder()
	router.ServeHTTP(resp10, req10)
	if resp10.Code != http.StatusOK {
		return fmt.Errorf("PATCH /users/:id failed: expected 200, got %d", resp10.Code)
	}
	log.Println("OK: PATCH /users/:id")

	// 6. DELETE /users/:id
	log.Println("Test: DELETE /users/:id")
	req11, _ := http.NewRequest("DELETE", "/users/"+createdUser.ID, nil)
	req11.Header.Set("Authorization", "Bearer "+token)
	resp11 := httptest.NewRecorder()
	router.ServeHTTP(resp11, req11)
	if resp11.Code != http.StatusOK {
		return fmt.Errorf("DELETE /users/:id failed: expected 200, got %d", resp11.Code)
	}
	log.Println("OK: DELETE /users/:id")

	return nil
}

// Setup testowego routera
func SetupTestRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/run-tests", RunTestsHandler)
	r.POST("/login", controllers.Login)

	albumRoutes := r.Group("/albums")
	albumRoutes.GET("", controllers.GetAlbums)
	albumRoutes.GET("/:id", controllers.GetAlbumByID)
	albumRoutes.Use(middleware.AuthMiddleware())
	{
		albumRoutes.POST("", middleware.RoleMiddleware("employee", "admin"), controllers.CreateAlbum)
		albumRoutes.PATCH("/:id", middleware.RoleMiddleware("employee", "admin"), controllers.UpdateAlbum)
		albumRoutes.DELETE("/:id", middleware.RoleMiddleware("employee", "admin"), controllers.DeleteAlbum)
	}

	userRoutes := r.Group("/users")
	userRoutes.Use(middleware.AuthMiddleware())
	{
		userRoutes.GET("", middleware.RoleMiddleware("admin"), controllers.GetUsers)
		userRoutes.GET("/:id", middleware.RoleMiddleware("admin"), controllers.GetUserByID)
		userRoutes.POST("", middleware.RoleMiddleware("admin"), controllers.CreateUser)
		userRoutes.PATCH("/:id", middleware.RoleMiddleware("admin"), controllers.UpdateUser)
		userRoutes.DELETE("/:id", middleware.RoleMiddleware("admin"), controllers.DeleteUser)
	}

	return r
}
