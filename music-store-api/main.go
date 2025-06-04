package main

import (
	"music-store-api/config"
	"music-store-api/controllers"
	_ "music-store-api/docs"
	"music-store-api/middleware"
	"music-store-api/tests"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Music Store REST API
// @version 1.0
// @description Music Store REST API to backendowy serwis RESTful do zarządzania zasobami internetowego sklepu muzycznego. Umożliwia zarządzanie albumami, użytkownikami, recenzjami, zamówieniami oraz procesami uwierzytelniania i autoryzacji. API zostało zaprojektowane do współpracy z frontendem aplikacji oraz systemami zewnętrznymi.
// @termsOfService http://example.com/terms/
// @contact.name Zespół Wsparcia Music Store
// @contact.url http://example.com/support
// @contact.email support@example.com
// @license.name MIT License
// @license.url https://opensource.org/licenses/MIT
// @host 193.28.226.78:25565
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Token JWT w formacie "Bearer <token>", wymagany do autoryzacji endpointów chronionych.
func main() {
	config.ConnectDB()
	defer config.DisconnectDB()

	controllers.InitAlbumCollection()
	controllers.InitUserCollection()
	controllers.InitOrderCollection()
	controllers.InitReviewCollection()

	r := gin.Default()

	r.Static("/static", "./static")
	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/login", controllers.Login)

	albumRoutes := r.Group("/albums")
	albumRoutes.GET("", controllers.GetAlbums)
	albumRoutes.GET("/:id", controllers.GetAlbumByID)
	albumRoutes.Use(middleware.AuthMiddleware())
	{
		albumRoutes.POST("", middleware.RoleMiddleware("employee", "admin"), controllers.CreateAlbum)
		albumRoutes.POST("/bulk", middleware.RoleMiddleware("employee", "admin"), controllers.CreateAlbumsBulk)
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

	orderRoutes := r.Group("/orders")
	orderRoutes.Use(middleware.AuthMiddleware())
	{
		orderRoutes.GET("/", middleware.RoleMiddleware("employee", "admin"), controllers.GetOrders)
		orderRoutes.GET("/:id", middleware.RoleMiddleware("employee", "admin"), controllers.GetOrderByID)
		orderRoutes.GET("/user/:userID", middleware.RoleMiddleware("employee", "admin"), controllers.GetOrdersByUserID)
		orderRoutes.POST("/", middleware.RoleMiddleware("customer", "employee", "admin"), controllers.CreateOrder)
		orderRoutes.PUT("/:id", middleware.RoleMiddleware("employee", "admin"), controllers.UpdateOrder)
		orderRoutes.DELETE("/:id", middleware.RoleMiddleware("admin"), controllers.DeleteOrder)
		orderRoutes.PATCH("/:id/status", middleware.RoleMiddleware("employee", "admin"), controllers.UpdateOrderStatus)
		orderRoutes.PUT("/:id/shipping", middleware.RoleMiddleware("employee", "admin"), controllers.UpdateOrderShipping)
	}

	reviewRoutes := r.Group("/reviews")
	reviewRoutes.GET("/album/:albumID", controllers.GetReviewsByAlbumID)
	reviewRoutes.GET("/user/:userID", controllers.GetReviewsByUserID)
	reviewRoutes.GET("/:id", controllers.GetReviewByID)
	reviewRoutes.GET("", controllers.GetReviews)
	reviewRoutes.Use(middleware.AuthMiddleware())
	{
		reviewRoutes.POST("", middleware.RoleMiddleware("customer", "employee", "admin"), controllers.CreateReview)
		reviewRoutes.PUT("/:id", middleware.RoleMiddleware("customer", "employee", "admin"), controllers.UpdateReview)
		reviewRoutes.DELETE("/:id", middleware.RoleMiddleware("employee", "admin"), controllers.DeleteReview)
	}

	dataRoutes := r.Group("/data")
	dataRoutes.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"))
	{
		dataRoutes.POST("/load", controllers.LoadTestData)
	}

	r.GET("/run-tests", tests.RunTestsHandler)

	r.Run(":25565")
}
