package controllers

import (
	"fmt"
	"go/gin-api/models"
	"log"
	"net/http"
	"net/smtp"
	"path"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const HOST, PORT = "smtp.gmail.com", 587
const SENDER = "PT. Digital Creative Studio <dirgantoro.facebook@gmail.com>"
const EMAIL, PASSWORD = "dirgantoro.facebook@gmail.com", "clzciwwmpbidehpk"

type CInput struct {
	Title  string `json:"title" binding:"required"`
	Author string `json:"author" binding:"required"`
}
type UInput struct {
	Title  string `json:"title"`
	Author string `json:"author"`
}
type SInput struct {
	Email string `json:"email" binding:"required"`
}

func Finds(c *gin.Context) {
	var b []models.Book
	models.DB.Find(&b)
	c.JSON(http.StatusOK, gin.H{"data": b})
}
func Find(c *gin.Context) {
	var b models.Book
	if err := models.DB.Where("id=?", c.Param("id")).First(&b).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": b})
}
func Create(c *gin.Context) {
	var i CInput
	if err := c.ShouldBindJSON(&i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	b := models.Book{Title: i.Title, Author: i.Author}
	models.DB.Create(&b)
	c.JSON(http.StatusOK, gin.H{"data": b})
}
func Update(c *gin.Context) {
	var b models.Book
	if err := models.DB.Where("id=?", c.Param("id")).First(&b).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var i UInput
	if err := c.ShouldBindJSON(&i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	models.DB.Model(&b).Updates(i)
	c.JSON(http.StatusOK, gin.H{"data": b})
}
func Delete(c *gin.Context) {
	var b models.Book
	if err := models.DB.Where("id=?", c.Param("id")).First(&b).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	models.DB.Delete(&b)
	c.JSON(http.StatusOK, gin.H{"data": true})
}
func SaveFileHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	extension := filepath.Ext(file.Filename)
	newFileName := uuid.New().String() + extension
	if err := c.SaveUploadedFile(file, path.Join("upload", newFileName)); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
func SendMail(c *gin.Context) {
	var input SInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	to, cc := []string{input.Email}, []string{}
	subject, message := "GoMailer", "Send mail with smtp golang\n\n"+
		"Best Regard,\nDirgantoro(CEO)\nRiko Primada(CTO)"
	if err := sendMailer(to, cc, subject, message); err != nil {
		log.Fatal(err.Error())
	}
	log.Println("Mail sent!")
	c.JSON(http.StatusOK, gin.H{"data": true})
}
func sendMailer(to []string, cc []string, subject, message string) error {
	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", HOST, PORT), smtp.PlainAuth("", EMAIL, PASSWORD, HOST), EMAIL,
		append(to, cc...), []byte(fmt.Sprintf("from: %s\nto: %s\ncc: %s\nsubject: %s\n\n%s", SENDER, strings.Join(to, ","), strings.Join(cc, ","), subject, message)))
	if err != nil {
		return err
	}
	return nil
}
