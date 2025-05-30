package main

import (
	"music-store-api/config"
	"music-store-api/controllers"
	_ "music-store-api/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	config.ConnectDB()
	defer config.DisconnectDB()

	controllers.InitAlbumCollection()
	controllers.InitUserCollection()
	controllers.InitCartCollection()
	controllers.InitOrderCollection()
	controllers.InitReviewCollection()

	r := gin.Default()

	r.Static("/static", "./static")
	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	albumRoutes := r.Group("/albums")
	{
		albumRoutes.GET("", controllers.GetAlbums)
		albumRoutes.GET("/:id", controllers.GetAlbumByID)
		albumRoutes.POST("", controllers.CreateAlbum)
		albumRoutes.POST("/bulk", controllers.CreateAlbumsBulk)
		albumRoutes.PATCH("/:id", controllers.UpdateAlbum)
		albumRoutes.DELETE("/:id", controllers.DeleteAlbum)
	}

	userRoutes := r.Group("/users")
	{
		userRoutes.GET("", controllers.GetUsers)
		userRoutes.GET("/:id", controllers.GetUserByID)
		userRoutes.POST("", controllers.CreateUser)
		userRoutes.PATCH("/:id", controllers.UpdateUser)
		userRoutes.DELETE("/:id", controllers.DeleteUser)
	}

	orderRoutes := r.Group("/orders")
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

	cartRoutes := r.Group("/carts")
	{
		cartRoutes.GET("", controllers.GetCarts)
		cartRoutes.GET("user/:userID", controllers.GetCartByUserID)
		cartRoutes.GET(":id", controllers.GetCartByID)
		cartRoutes.POST("", controllers.CreateCart)
		cartRoutes.PUT(":id", controllers.UpdateCart)
		cartRoutes.DELETE(":id", controllers.DeleteCart)
		cartRoutes.POST(":id/items", controllers.AddItemToCart)
		cartRoutes.DELETE(":id/items/:albumID", controllers.RemoveItemFromCart)
		cartRoutes.PUT(":id/items/:albumID", controllers.UpdateCartItemQuantity)
		cartRoutes.PUT(":id/total", controllers.UpdateCartTotal)
		cartRoutes.POST(":id/clear", controllers.ClearCart)
	}

	reviewRoutes := r.Group("/reviews")
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
	{
		dataRoutes.POST("/load", controllers.LoadTestData)
	}

	r.Run(":25565")
}
