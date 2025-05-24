package main

import (
	"music-store-api/config"
	"music-store-api/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDB()
	controllers.InitAlbumCollection()

	r := gin.Default()

	albumRoutes := r.Group("/albums")
	{
		albumRoutes.GET("", controllers.GetAlbums)
		albumRoutes.GET("/:id", controllers.GetAlbumByID)
		albumRoutes.POST("", controllers.CreateAlbum)
		albumRoutes.PUT("/:id", controllers.UpdateAlbum)
		albumRoutes.DELETE("/:id", controllers.DeleteAlbum)
	}
	r.Run(":8080")
}
