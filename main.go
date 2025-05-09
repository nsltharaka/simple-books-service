package main

import (
	"log"
	"log/slog"
	"net"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/nsltharaka/booksapi/database"
	"github.com/nsltharaka/booksapi/handlers"
	"github.com/nsltharaka/booksapi/services"
)

func envConfig() string {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file, using default values")
	}

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "3030"
	}

	return net.JoinHostPort(host, port)

}

func main() {

	serverAddr := envConfig()

	db, err := database.Connect()
	if err != nil {
		panic("error in creating database connection")
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

	app.Hooks().OnListen(func(listenData fiber.ListenData) error {
		logger.Info("Server started", slog.String("address", serverAddr))
		return nil
	})

	app.Listen(serverAddr)

}
