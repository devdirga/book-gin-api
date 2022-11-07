package models

import (
	"database/sql"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var SqliteDb *sql.DB

func Conn() error {
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		return err
	}
	SqliteDb = db
	return nil
}
