package controllers

import (
	f "fmt"
	m "go/gin-api/models"
	"log"
	h "net/http"
	"net/smtp"
	"path"
	"path/filepath"
	s "strings"

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
	var i CInput
	if e := c.ShouldBindJSON(&i); e != nil {
		c.JSON(h.StatusBadRequest, gin.H{"msg": e.Error()})
		return
	}
	trx, e := m.Db.Begin()
	if e != nil {
		c.JSON(h.StatusBadRequest, gin.H{"msg": e.Error()})
		return
	}
	stmnt, e := trx.Prepare(m.Insert)
	if e != nil {
		c.JSON(h.StatusBadRequest, gin.H{"msg": e.Error()})
		return
	}
	defer stmnt.Close()
	if _, e = stmnt.Exec(i.Title, i.Author); e != nil {
		c.JSON(h.StatusBadRequest, gin.H{"msg": e.Error()})
		return
	}
	trx.Commit()
	c.JSON(h.StatusOK, gin.H{"msg": m.MsgCreate})
}
func Finds(c *gin.Context) {
	rws, e := m.Db.Query(m.Finds)
	if e != nil {
		c.JSON(h.StatusBadRequest, gin.H{"msg": e.Error()})
		return
	}
	defer rws.Close()
	bks := make([]m.Book, 0)
	bk := m.Book{}
	for rws.Next() {
		rws.Scan(&bk.ID, &bk.Title, &bk.Author)
		bks = append(bks, bk)
	}
	c.JSON(h.StatusOK, gin.H{"msg": m.MsgFinds, "data": bks, "language": "世界"})
}
func Find(c *gin.Context) {
	var b m.Book
	row := m.Db.QueryRow(m.Find, c.Param("id"))
	if e := row.Scan(&b.ID, &b.Title, &b.Author); e != nil {
		c.JSON(h.StatusBadRequest, gin.H{"msg": e.Error()})
		return
	}
	c.JSON(h.StatusOK, gin.H{"msg": m.MsgFind, "data": b})
}
func Delete(c *gin.Context) {
	if _, e := m.Db.Exec(m.Delete, c.Param("id")); e != nil {
		c.AbortWithStatusJSON(h.StatusBadRequest, gin.H{"msg": e.Error()})
		return
	}
	c.JSON(h.StatusOK, gin.H{"msg": m.MsgDelete})
}
func Upload(c *gin.Context) {
	file, e := c.FormFile("file")
	if e != nil {
		c.AbortWithStatusJSON(h.StatusBadRequest, gin.H{"msg": e.Error()})
		return
	}
	ext := filepath.Ext(file.Filename)
	nFile := uuid.New().String() + ext
	if e := c.SaveUploadedFile(file, path.Join("upload", nFile)); e != nil {
		c.AbortWithStatusJSON(h.StatusBadRequest, gin.H{"msg": e.Error()})
		return
	}
	c.JSON(h.StatusOK, gin.H{"msg": m.MsgUpload})
}
func Mail(c *gin.Context) {
	var input SInput
	if e := c.ShouldBindJSON(&input); e != nil {
		c.JSON(h.StatusBadRequest, gin.H{"msg": e.Error()})
		return
	}
	to, cc := []string{input.Email}, []string{}
	if e := Mailer(to, cc, input.Subject, input.Message); e != nil {
		log.Fatal(e.Error())
	}
	c.JSON(h.StatusOK, gin.H{"msg": m.MsgMail})
}
func Mailer(to []string, cc []string, sbj, msg string) error {
	e := smtp.SendMail(
		f.Sprintf("%s:%d", m.Hst, m.Prt), smtp.PlainAuth("", m.Mail, m.Pwd, m.Hst), m.Mail,
		append(to, cc...), []byte(f.Sprintf("from: %s\nto: %s\ncc: %s\nsubject: %s\n\n%s", m.Sndr, s.Join(to, ","), s.Join(cc, ","), sbj, msg)))
	if e != nil {
		return e
	}
	return nil
}
