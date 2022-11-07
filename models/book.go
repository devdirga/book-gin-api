package models

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Book struct {
	ID     uint   `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var Create = "create table if not exists books (id integer primary key autoincrement, title text not null, author text not null);"
var Insert = "insert into books(title,author)VALUES(?,?)"
var Update = "update books set title=?,author=? where id=?"
var Delete = "delete from books where id=?"
var Finds = "select id,title,author from books"
var Find = "select id,title,author from books where id=?"

var Host = "smtp.gmail.com"
var Port = 587
var Sender = "PT. Digital Creative Studio <dirgantoro.facebook@gmail.com>"
var Email = "dirgantoro.facebook@gmail.com"
var Password = "clzciwwmpbidehpk"

func (b Book) HandleError(c *gin.Context, err error) {
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

func HandleErrorx(c *gin.Context, err error) {
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}
