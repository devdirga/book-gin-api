package models

type Book struct {
	ID     uint   `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

const (
	Create    = "create table if not exists books (id integer primary key autoincrement, title text not null, author text not null);"
	Insert    = "insert into books(title,author)values(?,?)"
	Update    = "update books set title=?,author=? where id=?"
	Delete    = "delete from books where id=?"
	Finds     = "select id,title,author from books"
	Find      = "select id,title,author from books where id=?"
	Hst       = "smtp.gmail.com"
	Prt       = 587
	Sndr      = "PT. Digital Creative Studio <dirgantoro.facebook@gmail.com>"
	Mail      = "dirgantoro.facebook@gmail.com"
	Pwd       = "clzciwwmpbidehpk"
	MsgCreate = "Success create book"
	MsgFinds  = "Success find books"
	MsgFind   = "Success find book"
	MsgDelete = "Success delete book"
	MsgUpload = "Success upload file"
	MsgMail   = "Success send mail"
)
