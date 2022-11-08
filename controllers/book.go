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

type CInput struct {
	Title  string `json:"title" binding:"required"`
	Author string `json:"author" binding:"required"`
}
type UInput struct {
	Title  string `json:"title"`
	Author string `json:"author"`
}
type SInput struct {
	Email   string `json:"email" binding:"required"`
	Message string `json:"message" binding:"required"`
}

func Insert(c *gin.Context) {
	var input CInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	transaction, err := models.Db.Begin()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	statement, err := transaction.Prepare(models.Insert)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	defer statement.Close()
	if _, err = statement.Exec(input.Title, input.Author); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	transaction.Commit()
	c.JSON(http.StatusOK, gin.H{"msg": models.MessageSuccessCreate})
}
func Finds(c *gin.Context) {
	rows, err := models.Db.Query(models.Finds)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	books := make([]models.Book, 0)
	book := models.Book{}
	for rows.Next() {
		rows.Scan(&book.ID, &book.Title, &book.Author)
		books = append(books, book)
	}
	c.JSON(http.StatusOK, gin.H{"msg": models.MessageSuccessFinds, "data": books})
}
func Find(c *gin.Context) {
	var book models.Book
	row := models.Db.QueryRow(models.Find, c.Param("id"))
	if err := row.Scan(&book.ID, &book.Title, &book.Author); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": models.MessageSuccessFind, "data": book})
}
func Delete(c *gin.Context) {
	if _, err := models.Db.Exec(models.Delete, c.Param("id")); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": models.MessageSuccessDelete})
}
func Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	ext := filepath.Ext(file.Filename)
	newFile := uuid.New().String() + ext
	if err := c.SaveUploadedFile(file, path.Join("upload", newFile)); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": models.MessageSuccessUpload})
}
func Mail(c *gin.Context) {
	var input SInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	to, cc := []string{input.Email}, []string{}
	sbj, msg := "GoMailer", input.Message+"\n\nBest Regard,\nDirgantoro\t(CEO)"
	if err := Mailer(to, cc, sbj, msg); err != nil {
		log.Fatal(err.Error())
	}
	c.JSON(http.StatusOK, gin.H{"msg": models.MessageSuccessMail})
}
func Mailer(to []string, cc []string, subject, message string) error {
	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", models.Host, models.Port), smtp.PlainAuth("", models.Email, models.Password, models.Host), models.Email,
		append(to, cc...), []byte(fmt.Sprintf("from: %s\nto: %s\ncc: %s\nsubject: %s\n\n%s", models.Sender, strings.Join(to, ","), strings.Join(cc, ","), subject, message)))
	if err != nil {
		return err
	}
	return nil
}
