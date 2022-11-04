package main

import (
	"go/gin-api/controllers"
	"go/gin-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	models.ConnectDatabase()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  200,
			"message": "pong",
		})
	})
	v1 := r.Group("/v1")
	{
		v1.GET("books", controllers.Finds)
		v1.GET("book/:id", controllers.Find)
		v1.POST("book", controllers.Create)
		v1.PUT("book/:id", controllers.Update)
		v1.DELETE("book/:id", controllers.Delete)
		v1.POST("upload", controllers.SaveFileHandler)
		v1.POST("email", controllers.SendMail)
	}
	r.Run()
}
