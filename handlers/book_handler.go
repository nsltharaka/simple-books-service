package handlers

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/nsltharaka/booksapi/models"
	"github.com/nsltharaka/booksapi/services"
)

func ErrorHandler(c *fiber.Ctx, err error) error {

	code := fiber.StatusInternalServerError

	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)

	return c.Status(code).JSON(apiResponse{
		Message: "error",
		Error:   err.Error(),
	})
}

type BookHandler struct {
	bookService services.IBookService
	validate    *validator.Validate
}

func NewBookHandler(service services.IBookService, validator *validator.Validate) *BookHandler {
	return &BookHandler{
		bookService: service,
		validate:    validator,
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

	page := c.QueryInt("page")
	if page <= 0 {
		page = 1
	}

	limit := c.QueryInt("limit")
	switch {
	case limit > 100:
		limit = 100
	case limit <= 0:
		limit = 10
	}

	books, err := handler.bookService.GetAllBooks(page, limit)
	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(apiResponse{
		Message: "success",
		Data:    books,
	})
}

func (handler *BookHandler) getBook(c *fiber.Ctx) error {
	bookId, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid parameter")
	}

	book, err := handler.bookService.GetBook(uint(bookId))
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return err
	}

	return c.Status(http.StatusOK).JSON(apiResponse{
		Message: "success",
		Data:    book,
	})
}

func (handler *BookHandler) newBook(c *fiber.Ctx) error {
	var book models.Book
	if err := c.BodyParser(&book); err != nil {
		if errors.Is(err, fiber.ErrUnprocessableEntity) {
			return fiber.NewError(fiber.StatusBadRequest, "invalid Content-Type")
		}
	}

	err := handler.validate.Struct(&book)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid payload")
	}

	createdBook, err := handler.bookService.CreateBook(&book)
	if err != nil {
		return err
	}

	return c.Status(http.StatusCreated).JSON(apiResponse{
		Message: "success",
		Data:    createdBook,
	})

}

func (handler *BookHandler) updateBook(c *fiber.Ctx) error {
	bookId, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid parameter")
	}

	var book models.Book
	if err := c.BodyParser(&book); err != nil {
		if errors.Is(err, fiber.ErrUnprocessableEntity) {
			return fiber.NewError(fiber.StatusBadRequest, "invalid Content-Type")
		}
	}

	err = handler.validate.Struct(&book)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid payload")
	}

	book.ID = uint(bookId)
	updatedBook, err := handler.bookService.UpdateBook(&book)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return err
	}

	return c.Status(fiber.StatusOK).JSON(apiResponse{
		Message: "success",
		Data:    updatedBook,
	})

}

func (handler *BookHandler) deleteBook(c *fiber.Ctx) error {
	bookId, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid parameter")
	}

	book, err := handler.bookService.DeleteBook(uint(bookId))
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return err
	}

	return c.Status(http.StatusOK).JSON(apiResponse{
		Message: "success",
		Data:    book,
	})
}

type apiResponse struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
	Data    any    `json:"data,omitempty"`
}
