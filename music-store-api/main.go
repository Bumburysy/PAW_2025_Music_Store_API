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
	controllers.InitAlbumCollection()
	controllers.InitUserCollection()

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

	dataRoutes := r.Group("/data")
	{
		dataRoutes.POST("/load", controllers.LoadTestData)
	}

	r.Run(":25565")
}
