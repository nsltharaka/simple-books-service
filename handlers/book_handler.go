package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/nsltharaka/booksapi/models"
	"github.com/nsltharaka/booksapi/services"
)

type BookHandler struct {
	bookService *services.BookService
}

func NewBookHandler(service *services.BookService) *BookHandler {
	return &BookHandler{
		bookService: service,
	}
}

func (handler *BookHandler) SetupRoutes(router fiber.Router) {
	router.Get("/books", handler.getAllBooks)
	router.Get("/books/:id", handler.getBook)
	router.Post("/books", handler.newBook)
	router.Put("/books/:id", handler.updateBook)
	router.Delete("/books/:id", handler.deleteBook)
}

func (handler *BookHandler) getAllBooks(c *fiber.Ctx) error {
	books, err := handler.bookService.GetAllBooks()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err.Error())
	}

	return c.Status(http.StatusOK).JSON(books)
}

func (handler *BookHandler) getBook(c *fiber.Ctx) error {
	param := c.Params("id", "")
	bookId, err := strconv.Atoi(param)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(err.Error())
	}

	book, err := handler.bookService.GetBook(uint(bookId))
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			return c.Status(http.StatusNotFound).JSON(err.Error())
		}
		return c.Status(http.StatusInternalServerError).JSON(err.Error())
	}

	return c.Status(http.StatusOK).JSON(book)
}

func (handler *BookHandler) newBook(c *fiber.Ctx) error {
	var book models.Book
	if err := c.BodyParser(&book); err != nil {
		return c.Status(http.StatusBadRequest).JSON(err.Error())
	}

	createdBook, err := handler.bookService.CreateBook(&models.Book{
		Title:  book.Title,
		Author: book.Author,
		Year:   book.Year,
	})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err.Error())
	}

	return c.Status(http.StatusCreated).JSON(createdBook)

}

func (handler *BookHandler) updateBook(c *fiber.Ctx) error {
	return nil
}

func (handler *BookHandler) deleteBook(c *fiber.Ctx) error {
	return nil
}
