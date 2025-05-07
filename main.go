package main

import (
	"log/slog"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/nsltharaka/booksapi/database"
	"github.com/nsltharaka/booksapi/handlers"
	"github.com/nsltharaka/booksapi/models"
	"github.com/nsltharaka/booksapi/services"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		panic("error passing environment variables")
	}

	db, err := database.Connection()
	if err != nil {
		slog.Error("error creating database connection", "err", err)
		return
	}

	// run database migrations
	db.AutoMigrate(&models.Book{})

	app := fiber.New()
	apiV1 := app.Group("/api").Group("/v1")

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	validator := validator.New(validator.WithRequiredStructEnabled())

	bookService := services.NewBookService(db, logger)
	bookHandler := handlers.NewBookHandler(bookService, validator)

	bookHandler.SetupRoutes(apiV1)

	app.Listen(":3000")

}
