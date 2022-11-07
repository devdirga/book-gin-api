package main

import (
	"go/gin-api/controllers"
	"go/gin-api/models"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	models.Conn()
	v1 := r.Group("/v1")
	{
		v1.POST("book", controllers.Insert)
		v1.GET("books", controllers.Finds)
		v1.GET("book/:id", controllers.Find)
		v1.DELETE("book/:id", controllers.Delete)
		v1.POST("upload", controllers.Upload)
		v1.POST("email", controllers.Mail)
	}
	v2 := r.Group("/v2")
	{
		v2.POST("article", controllers.Insert)
		v2.GET("articles", controllers.Finds)
		v2.GET("article/:id", controllers.Find)
		v2.DELETE("article/:id", controllers.Delete)
	}
	r.Run()
}
