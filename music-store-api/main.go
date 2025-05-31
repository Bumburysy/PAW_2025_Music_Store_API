package main

import (
	"music-store-api/config"
	"music-store-api/controllers"
	_ "music-store-api/docs"
	"music-store-api/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Music Store API
// @version 1.0
// @description API do zarządzania sklepem muzycznym

// @host localhost:25565

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	config.ConnectDB()
	defer config.DisconnectDB()

	controllers.InitAlbumCollection()
	controllers.InitUserCollection()
	controllers.InitOrderCollection()
	controllers.InitReviewCollection()

	r := gin.Default()

	r.POST("/login", controllers.Login)

	r.Static("/static", "./static")
	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	albumRoutes := r.Group("/albums")
	albumRoutes.GET("", controllers.GetAlbums)
	albumRoutes.GET("/:id", controllers.GetAlbumByID)

	albumRoutes.Use(middleware.AuthMiddleware())
	{
		albumRoutes.POST("", controllers.CreateAlbum)
		albumRoutes.POST("/bulk", controllers.CreateAlbumsBulk)
		albumRoutes.PATCH("/:id", controllers.UpdateAlbum)
		albumRoutes.DELETE("/:id", controllers.DeleteAlbum)
	}

	userRoutes := r.Group("/users")
	userRoutes.Use(middleware.AuthMiddleware())
	{
		userRoutes.GET("", controllers.GetUsers)
		userRoutes.GET("/:id", controllers.GetUserByID)
		userRoutes.POST("", controllers.CreateUser)
		userRoutes.PATCH("/:id", controllers.UpdateUser)
		userRoutes.DELETE("/:id", controllers.DeleteUser)
	}

	orderRoutes := r.Group("/orders")
	orderRoutes.Use(middleware.AuthMiddleware())
	{
		orderRoutes.GET("/", controllers.GetOrders)
		orderRoutes.GET("/:id", controllers.GetOrderByID)
		orderRoutes.GET("/user/:userID", controllers.GetOrdersByUserID)
		orderRoutes.POST("/", controllers.CreateOrder)
		orderRoutes.PUT("/:id", controllers.UpdateOrder)
		orderRoutes.DELETE("/:id", controllers.DeleteOrder)
		orderRoutes.PATCH("/:id/status", controllers.UpdateOrderStatus)
		orderRoutes.PUT("/:id/shipping", controllers.UpdateOrderShipping)
	}

	reviewRoutes := r.Group("/reviews")
	reviewRoutes.Use(middleware.AuthMiddleware())
	{
		reviewRoutes.GET("/album/:albumID", controllers.GetReviewsByAlbumID)
		reviewRoutes.GET("/user/:userID", controllers.GetReviewsByUserID)
		reviewRoutes.GET("/:id", controllers.GetReviewByID)
		reviewRoutes.GET("", controllers.GetReviews)
		reviewRoutes.POST("", controllers.CreateReview)
		reviewRoutes.PUT("/:id", controllers.UpdateReview)
		reviewRoutes.DELETE("/:id", controllers.DeleteReview)
	}

	dataRoutes := r.Group("/data")
	//dataRoutes.Use(middleware.AuthMiddleware())
	{
		dataRoutes.POST("/load", controllers.LoadTestData)
	}

	r.Run(":25565")
}
