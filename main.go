package main

import (
	"log/slog"

	"github.com/nsltharaka/booksapi/database"
	"github.com/nsltharaka/booksapi/models"
)

func main() {
	db, err := database.Connection()
	if err != nil {
		slog.Error("error creating database connection", "err", err)
		return
	}

	// run migrations
	db.AutoMigrate(&models.Book{})
}
