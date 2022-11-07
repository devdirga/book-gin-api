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
		v1.POST("upload", controllers.SaveFileHandler)
		v1.POST("email", controllers.SendMail)
	}
	r.Run()
}
