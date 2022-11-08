package models

type Book struct {
	ID     uint   `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var Create = "create table if not exists books (id integer primary key autoincrement, title text not null, author text not null);"
var Insert = "insert into books(title,author)values(?,?)"
var Update = "update books set title=?,author=? where id=?"
var Delete = "delete from books where id=?"
var Finds = "select id,title,author from books"
var Find = "select id,title,author from books where id=?"

var Hst = "smtp.gmail.com"
var Prt = 587
var Sndr = "PT. Digital Creative Studio <dirgantoro.facebook@gmail.com>"
var Mail = "dirgantoro.facebook@gmail.com"
var Pwd = "clzciwwmpbidehpk"

var MsgCreate = "Success create book"
var MsgFinds = "Success find books"
var MsgFind = "Success find book"
var MsgDelete = "Success delete book"
var MsgUpload = "Success upload file"
var MsgMail = "Success send mail"
