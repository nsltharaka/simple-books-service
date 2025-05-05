package database

import (
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Connection() (*gorm.DB, error) {
	dbFile := os.Getenv("SQLITE_FILENAME")
	if dbFile == "" {
		return nil, fmt.Errorf("environment variable SQLITE_FILENAME is not set")
	}

	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
