package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/nsltharaka/booksapi/database"
	"github.com/nsltharaka/booksapi/models"
	"github.com/nsltharaka/booksapi/services"
	"github.com/stretchr/testify/assert"
)

func setupTestService(t *testing.T) (*services.BookService, func()) {
	filename := "test.db"
	os.Setenv("SQLITE_FILENAME", filename)
	db, _ := database.Connection()
	db.AutoMigrate(&models.Book{})

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	service := services.NewBookService(db, logger)

	books := []*models.Book{
		{Title: "Book One", Author: "Author A", Year: 2021},
		{Title: "Book Two", Author: "Author B", Year: 2022},
		{Title: "Book Three", Author: "Author C", Year: 2023},
	}

	// Create books
	for _, book := range books {
		service.CreateBook(book)
	}

	cleanup := func() {
		os.Unsetenv("SQLITE_FILENAME")
		os.Remove(filename)
	}

	return service, cleanup
}

func setupTestApp(t *testing.T) *fiber.App {
	service, cleanup := setupTestService(t)
	t.Cleanup(cleanup)

	validator := validator.New(validator.WithRequiredStructEnabled())

	app := fiber.New()
	handler := NewBookHandler(service, validator)
	handler.SetupRoutes(app)

	return app
}

func TestGetBooksHandler(t *testing.T) {

	t.Run("Get all books", func(t *testing.T) {
		app := setupTestApp(t)
		req := httptest.NewRequest("GET", "/books", nil)
		res, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("Get existing book", func(t *testing.T) {

		expected := models.Book{ID: 1, Title: "Book One", Author: "Author A", Year: 2021}

		app := setupTestApp(t)
		req := httptest.NewRequest("GET", "/books/1", nil)

		res, err := app.Test(req, -1)
		assert.NoError(t, err)

		body, err := io.ReadAll(res.Body)
		assert.NoError(t, err)

		var result models.Book
		json.Unmarshal(body, &result)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, result, expected)

	})
	t.Run("Get non-existing book", func(t *testing.T) {

		app := setupTestApp(t)
		req := httptest.NewRequest("GET", "/books/99", nil)

		res, err := app.Test(req, -1)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusNotFound, res.StatusCode)

	})

}

func TestCreateBook(t *testing.T) {

	book := struct {
		Title  string
		Author string
		Year   int
	}{
		Title:  "Book Four",
		Author: "Author D",
		Year:   2024,
	}

	jsonString, _ := json.Marshal(&book)

	app := setupTestApp(t)

	t.Run("creating a book without request payload", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/books", nil)
		req.Header.Add("Content-Type", "application/json")

		res, err := app.Test(req, -1)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})

	t.Run("creating a book with malformed request payload", func(t *testing.T) {
		testPayloads := []string{
			`{"malformedTitle" : "test", "malformedAuthor" : "test", "malformedYear" : "test"}`, // malformed fields
			`{"title" : "", "author" : "", "year" : ""}`,                                        // with zero values
			`{"title" : "", "author" : "testAuthor", "year" : 2025}`,                            // one field with zero value
			`{"title" : " ", "author" : "testAuthor", "year" : 2025}`,                           // string field with space
			`{"title" : "testTitle", "author" : "testAuthor", "year" : 0}`,                      // number field with zero value
			`{"title": "testTitle", "author": "testAuthor", "year":}`,                           // missing value
			`{"title": "testTitle", "author": "testAuthor", "year": 2025`,                       // missing closing brace
			`"title": "testTitle", "author": "testAuthor", "year": 2025}`,                       // missing opening brace
			`{title: "testTitle", author: "testAuthor", year: 2025}`,                            // keys not in quotes
			`{}`, // empty object
		}

		for idx, malformedJson := range testPayloads {
			idx := idx
			malformedJson := malformedJson

			t.Run(fmt.Sprintf("test payload : %d", idx+1), func(t *testing.T) {
				req := httptest.NewRequest("POST", "/books", strings.NewReader(malformedJson))
				req.Header.Set("Content-Type", "application/json")

				res, err := app.Test(req, -1)

				assert.NoError(t, err)
				assert.Equal(t, http.StatusBadRequest, res.StatusCode)
			})
		}
	})

	t.Run("creating a book with correct payload", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/books", bytes.NewReader(jsonString))
		req.Header.Add("Content-Type", "application/json")

		res, err := app.Test(req, -1)
		assert.NoError(t, err)

		resBody, _ := io.ReadAll(res.Body)

		var result models.Book
		json.Unmarshal(resBody, &result)

		assert.Equal(t, http.StatusCreated, res.StatusCode)
		assert.Equal(t, uint(4), result.ID)
		assert.Equal(t, book.Title, result.Title)
		assert.Equal(t, book.Author, result.Author)
		assert.Equal(t, book.Year, result.Year)
	})

}

func TestUpdateBook(t *testing.T) {

	updatedBook := struct {
		Title  string
		Author string
		Year   int
	}{
		Title:  "Updated Book One",
		Author: "Updated Author A",
		Year:   2030,
	}

	body, _ := json.Marshal(&updatedBook)

	t.Run("update request with invalid param", func(t *testing.T) {
		app := setupTestApp(t)

		req := httptest.NewRequest("PUT", "/books/xx", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req, -1)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})

	t.Run("updating an existing book without request body", func(t *testing.T) {
		app := setupTestApp(t)

		req := httptest.NewRequest("PUT", "/books/1", nil)
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})

	t.Run("updating an existing book with malformed request body", func(t *testing.T) {
		app := setupTestApp(t)

		testPayloads := []string{
			`{"malformedTitle" : "test", "malformedAuthor" : "test", "malformedYear" : "test"}`, // malformed fields
			`{"title" : "", "author" : "", "year" : ""}`,                                        // with zero values
			`{"title" : "", "author" : "testAuthor", "year" : 2025}`,                            // one field with zero value
			`{"title" : " ", "author" : "testAuthor", "year" : 2025}`,                           // string field with space
			`{"title" : "testTitle", "author" : "testAuthor", "year" : 0}`,                      // number field with zero value
			`{"title": "testTitle", "author": "testAuthor", "year":}`,                           // missing value
			`{"title": "testTitle", "author": "testAuthor", "year": 2025`,                       // missing closing brace
			`"title": "testTitle", "author": "testAuthor", "year": 2025}`,                       // missing opening brace
			`{title: "testTitle", author: "testAuthor", year: 2025}`,                            // keys not in quotes
			`{}`, // empty object
		}

		for idx, malformedJson := range testPayloads {
			idx := idx
			malformedJson := malformedJson

			t.Run(fmt.Sprintf("test payload : %d", idx+1), func(t *testing.T) {
				req := httptest.NewRequest("PUT", "/books/1", strings.NewReader(malformedJson))
				req.Header.Set("Content-Type", "application/json")

				res, err := app.Test(req, -1)

				assert.NoError(t, err)
				assert.Equal(t, http.StatusBadRequest, res.StatusCode)
			})
		}

	})

	t.Run("updating an existing book", func(t *testing.T) {
		app := setupTestApp(t)

		req := httptest.NewRequest("PUT", "/books/1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		var result models.Book
		resBody, _ := io.ReadAll(res.Body)
		json.Unmarshal(resBody, &result)

		assert.Equal(t, uint(1), result.ID)
		assert.Equal(t, updatedBook.Title, result.Title)
		assert.Equal(t, updatedBook.Author, result.Author)
		assert.Equal(t, updatedBook.Year, result.Year)
	})

}

func TestDeleteBook(t *testing.T) {
	app := setupTestApp(t)

	t.Run("Deleting a non-existing book", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/books/99", nil)
		res, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, res.StatusCode)
	})

	t.Run("Deleting an existing book", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/books/1", nil)
		res, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, res.StatusCode)

		getReq := httptest.NewRequest("GET", "/books/1", nil)
		getRes, err := app.Test(getReq, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, getRes.StatusCode)
	})

}
