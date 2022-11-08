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
	if e := c.ShouldBindJSON(&input); e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": e.Error()})
		return
	}
	transaction, e := m.Db.Begin()
	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": e.Error()})
		return
	}
	stmnt, e := transaction.Prepare(m.Insert)
	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": e.Error()})
		return
	}
	defer stmnt.Close()
	if _, e = stmnt.Exec(input.Title, input.Author); e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": e.Error()})
		return
	}
	transaction.Commit()
	c.JSON(http.StatusOK, gin.H{"msg": m.MsgCreate})
}
func Finds(c *gin.Context) {
	rows, e := m.Db.Query(m.Finds)
	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": e.Error()})
		return
	}
	defer rows.Close()
	books := make([]m.Book, 0)
	book := m.Book{}
	for rows.Next() {
		rows.Scan(&book.ID, &book.Title, &book.Author)
		books = append(books, book)
	}
	c.JSON(http.StatusOK, gin.H{"msg": m.MsgFinds, "data": books, "language": "世界"})
}
func Find(c *gin.Context) {
	var book m.Book
	row := m.Db.QueryRow(m.Find, c.Param("id"))
	if e := row.Scan(&book.ID, &book.Title, &book.Author); e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": e.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": m.MsgFind, "data": book})
}
func Delete(c *gin.Context) {
	if _, e := m.Db.Exec(m.Delete, c.Param("id")); e != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": e.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": m.MsgDelete})
}
func Upload(c *gin.Context) {
	file, e := c.FormFile("file")
	if e != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": e.Error()})
		return
	}
	ext := filepath.Ext(file.Filename)
	nFile := uuid.New().String() + ext
	if e := c.SaveUploadedFile(file, path.Join("upload", nFile)); e != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": e.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": m.MsgUpload})
}
func Mail(c *gin.Context) {
	var input SInput
	if e := c.ShouldBindJSON(&input); e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": e.Error()})
		return
	}
	to, cc := []string{input.Email}, []string{}
	if e := Mailer(to, cc, input.Subject, input.Message); e != nil {
		log.Fatal(e.Error())
	}
	c.JSON(http.StatusOK, gin.H{"msg": m.MsgMail})
}
func Mailer(to []string, cc []string, sbj, msg string) error {
	e := smtp.SendMail(
		f.Sprintf("%s:%d", m.Hst, m.Prt), smtp.PlainAuth("", m.Mail, m.Pwd, m.Hst), m.Mail,
		append(to, cc...), []byte(f.Sprintf("from: %s\nto: %s\ncc: %s\nsubject: %s\n\n%s", m.Sndr, strings.Join(to, ","), strings.Join(cc, ","), sbj, msg)))
	if e != nil {
		return e
	}
	return nil
}
