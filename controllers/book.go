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

const SMTP_HOST, SMTP_PORT = "smtp.gmail.com", 587
const SENDER_NAME = "PT. Digital Creative Studio <dirgantoro.facebook@gmail.com>"
const AUTH_EMAIL, AUTH_PASSWORD = "dirgantoro.facebook@gmail.com", "clzciwwmpbidehpk"

type CreateBookInput struct {
	Title  string `json:"title" binding:"required"`
	Author string `json:"author" binding:"required"`
}
type UpdateBookInput struct {
	Title  string `json:"title"`
	Author string `json:"author"`
}
type SendMailInput struct {
	Email string `json:"email" binding:"required"`
}

func Finds(c *gin.Context) {
	var books []models.Book
	models.DB.Find(&books)
	c.JSON(http.StatusOK, gin.H{"data": books})
}
func Find(c *gin.Context) {
	var book models.Book
	if err := models.DB.Where("id=?", c.Param("id")).First(&book).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": book})
}
func Create(c *gin.Context) {
	var input CreateBookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	book := models.Book{Title: input.Title, Author: input.Author}
	models.DB.Create(&book)
	c.JSON(http.StatusOK, gin.H{"data": book})
}
func Update(c *gin.Context) {
	var book models.Book
	if err := models.DB.Where("id=?", c.Param("id")).First(&book).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var input UpdateBookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	models.DB.Model(&book).Updates(input)
	c.JSON(http.StatusOK, gin.H{"data": book})
}
func Delete(c *gin.Context) {
	var book models.Book
	if err := models.DB.Where("id=?", c.Param("id")).First(&book).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	models.DB.Delete(&book)
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
	var input SendMailInput
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
		fmt.Sprintf("%s:%d", SMTP_HOST, SMTP_PORT),
		smtp.PlainAuth("", AUTH_EMAIL, AUTH_PASSWORD, SMTP_HOST),
		AUTH_EMAIL,
		append(to, cc...), []byte(fmt.Sprintf("from: %s\nto: %s\ncc: %s\nsubject: %s\n\n%s", SENDER_NAME, strings.Join(to, ","), strings.Join(cc, ","), subject, message)))
	if err != nil {
		return err
	}
	return nil
}
