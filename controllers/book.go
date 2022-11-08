package controllers

import (
	f "fmt"
	m "go/gin-api/models"
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
	Subject string `json:"subject" binding:"required"`
	Message string `json:"message" binding:"required"`
}

func Insert(c *gin.Context) {
	var input CInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	transaction, err := m.Db.Begin()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	stmnt, err := transaction.Prepare(m.Insert)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	defer stmnt.Close()
	if _, err = stmnt.Exec(input.Title, input.Author); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	transaction.Commit()
	c.JSON(http.StatusOK, gin.H{"msg": m.MsgSccCreate})
}
func Finds(c *gin.Context) {
	rows, err := m.Db.Query(m.Finds)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	defer rows.Close()
	books := make([]m.Book, 0)
	book := m.Book{}
	for rows.Next() {
		rows.Scan(&book.ID, &book.Title, &book.Author)
		books = append(books, book)
	}
	c.JSON(http.StatusOK, gin.H{"msg": m.MsgSccFinds, "data": books, "language": "世界"})
}
func Find(c *gin.Context) {
	var book m.Book
	row := m.Db.QueryRow(m.Find, c.Param("id"))
	if err := row.Scan(&book.ID, &book.Title, &book.Author); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": m.MsgSccFind, "data": book})
}
func Delete(c *gin.Context) {
	if _, err := m.Db.Exec(m.Delete, c.Param("id")); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": m.MsgSccDelete})
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
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": m.MsgSccUpload})
}
func Mail(c *gin.Context) {
	var input SInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	to, cc := []string{input.Email}, []string{}
	if err := Mailer(to, cc, input.Subject, input.Message); err != nil {
		log.Fatal(err.Error())
	}
	c.JSON(http.StatusOK, gin.H{"msg": m.MsgSccMail})
}
func Mailer(to []string, cc []string, subject, message string) error {
	err := smtp.SendMail(
		f.Sprintf("%s:%d", m.Host, m.Port), smtp.PlainAuth("", m.Email, m.Password, m.Host), m.Email,
		append(to, cc...), []byte(f.Sprintf("from: %s\nto: %s\ncc: %s\nsubject: %s\n\n%s", m.Sender, strings.Join(to, ","), strings.Join(cc, ","), subject, message)))
	if err != nil {
		return err
	}
	return nil
}
