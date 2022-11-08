package models

import (
	"database/sql"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var Db *sql.DB

func Conn() error {
	d, e := sql.Open("sqlite3", "book.db")
	if e != nil {
		return e
	}
	t, e := d.Begin()
	if e != nil {
		return e
	}
	s, e := t.Prepare(Create)
	if e != nil {
		return e
	}
	defer s.Close()
	if _, e = s.Exec(); e != nil {
		return e
	}
	t.Commit()
	Db = d
	return nil
}
