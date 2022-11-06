package models

import (
	"database/sql"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var DB *gorm.DB
var SqliteDb *sql.DB

func Connect() {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("Failed to connect to database!")
	}
	DB = db.AutoMigrate(&Book{})
}
func ConnectDB() error {
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		return err
	}
	SqliteDb = db
	return nil
}
