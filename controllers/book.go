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
	Email string `json:"email" binding:"required"`
}

// StartTLS sends the STARTTLS command and encrypts all further communication.
// Only servers that advertise the STARTTLS extension support this function.
func Insert(c *gin.Context) {
	var input CInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	transaction, err := models.SqliteDb.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	statement, err := transaction.Prepare(models.Insert)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer statement.Close()
	_, err = statement.Exec(input.Title, input.Author)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	transaction.Commit()
	c.JSON(http.StatusOK, gin.H{"data": true})
}

// StartTLS sends the STARTTLS command and encrypts all further communication.
// Only servers that advertise the STARTTLS extension support this function.
func Finds(c *gin.Context) {
	rows, err := models.SqliteDb.Query(models.Finds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	books := make([]models.Book, 0)
	book := models.Book{}
	for rows.Next() {
		rows.Scan(&book.ID, &book.Title, &book.Author)
		books = append(books, book)
	}
	c.JSON(http.StatusOK, gin.H{"data": books})
}

// StartTLS sends the STARTTLS command and encrypts all further communication.
// Only servers that advertise the STARTTLS extension support this function.
func Find(c *gin.Context) {
	var book models.Book
	row := models.SqliteDb.QueryRow(models.Find, c.Param("id"))
	err := row.Scan(&book.ID, &book.Title, &book.Author)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": book})
}

// StartTLS sends the STARTTLS command and encrypts all further communication.
// Only servers that advertise the STARTTLS extension support this function.
func Delete(c *gin.Context) {
	if _, err := models.SqliteDb.Exec(models.Delete, c.Param("id")); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "success"})
}

// StartTLS sends the STARTTLS command and encrypts all further communication.
// Only servers that advertise the STARTTLS extension support this function.
func SaveFileHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ext := filepath.Ext(file.Filename)
	newFile := uuid.New().String() + ext
	if err := c.SaveUploadedFile(file, path.Join("upload", newFile)); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "ok"})
}

// StartTLS sends the STARTTLS command and encrypts all further communication.
// Only servers that advertise the STARTTLS extension support this function.
func SendMail(c *gin.Context) {
	var input SInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	to, cc := []string{input.Email}, []string{}
	sbj, msg := "GoMailer", "Send mail with smtp golang\n\nBest Regard,\nDirgantoro(CEO)\nRiko Primada(CTO)"
	if err := sendMailer(to, cc, sbj, msg); err != nil {
		log.Fatal(err.Error())
	}
	c.JSON(http.StatusOK, gin.H{"msg": "success"})
}

// StartTLS sends the STARTTLS command and encrypts all further communication.
// Only servers that advertise the STARTTLS extension support this function.
func sendMailer(to []string, cc []string, subject, message string) error {
	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", models.Host, models.Port), smtp.PlainAuth("", models.Email, models.Password, models.Host), models.Email,
		append(to, cc...), []byte(fmt.Sprintf("from: %s\nto: %s\ncc: %s\nsubject: %s\n\n%s", models.Sender, strings.Join(to, ","), strings.Join(cc, ","), subject, message)))
	if err != nil {
		return err
	}
	return nil
}
