package main

import (
	"log/slog"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/nsltharaka/booksapi/database"
	"github.com/nsltharaka/booksapi/handlers"
	"github.com/nsltharaka/booksapi/services"
)

func main() {

	if err := godotenv.Load(); err != nil {
		panic("error passing environment variables")
	}

	db, err := database.Connect()
	if err != nil {
		panic("error passing environment variables")
	}

	config := fiber.Config{
		ErrorHandler: handlers.ErrorHandler,
	}

	app := fiber.New(config)
	app.Use(cors.New())

	apiV1 := app.Group("/api").Group("/v1")

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	validator := validator.New(validator.WithRequiredStructEnabled())

	bookService := services.NewBookService(db, logger)
	bookHandler := handlers.NewBookHandler(bookService, validator)
	bookHandler.SetupRoutes(apiV1)

	app.Listen(":3000")

}
