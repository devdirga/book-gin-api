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
	r.GET("books", controllers.FindBooks)
	r.GET("book/:id", controllers.FindBook)
	r.POST("book", controllers.CreateBook)
	r.PUT("book/:id", controllers.UpdateBook)
	r.DELETE("book/:id", controllers.DeleteBook)
	r.Run()
}
