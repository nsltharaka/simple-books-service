package database

import (
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nsltharaka/booksapi/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	dbFile := os.Getenv("SQLITE_FILENAME")
	if dbFile == "" {
		dbFile = "books.db"
	}

	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&models.Book{})

	return db, nil
}
