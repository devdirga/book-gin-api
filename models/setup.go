package models

import (
	"database/sql"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var Db *sql.DB

func Conn() error {
	db, err := sql.Open("sqlite3", "book.db")
	if err != nil {
		return err
	}
	trans, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := trans.Prepare(Create)
	if err != nil {
		return err
	}
	defer stmt.Close()
	if _, err = stmt.Exec(); err != nil {
		return err
	}
	trans.Commit()
	Db = db
	return nil
}
