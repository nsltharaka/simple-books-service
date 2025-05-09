package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/nsltharaka/booksapi/models"
	"github.com/nsltharaka/booksapi/services"
	"github.com/stretchr/testify/assert"
)

func setupTestApp(t *testing.T) *fiber.App {
	validator := validator.New(validator.WithRequiredStructEnabled())
	mockedBookService := NewMockedBookService()

	handler := NewBookHandler(mockedBookService, validator)

	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler,
	})

	handler.SetupRoutes(app)
	return app
}

func TestGetBooksHandler(t *testing.T) {

	t.Run("Get all books", func(t *testing.T) {
		app := setupTestApp(t)
		req := httptest.NewRequest("GET", "/books", nil)
		res, err := app.Test(req, -1)

		var apiResponse apiResponse
		json.NewDecoder(res.Body).Decode(&apiResponse)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "success", apiResponse.Message)
		assert.NotEmpty(t, apiResponse.Data)
		assert.Len(t, apiResponse.Data, 3)
	})

	t.Run("Get all books with pagination", func(t *testing.T) {
		app := setupTestApp(t)
		req := httptest.NewRequest("GET", "/books?page=1&limit=2", nil)
		res, err := app.Test(req, -1)

		var apiResponse apiResponse
		json.NewDecoder(res.Body).Decode(&apiResponse)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "success", apiResponse.Message)
		assert.NotEmpty(t, apiResponse.Data)
		assert.Len(t, apiResponse.Data, 2)
	})

	t.Run("Get existing book", func(t *testing.T) {

		expected := models.Book{Title: "Book One", Author: "Author A", Year: 2021}

		app := setupTestApp(t)
		req := httptest.NewRequest("GET", "/books/1", nil)

		res, err := app.Test(req, -1)
		assert.NoError(t, err)

		var apiResponse struct {
			Data    models.Book `json:"data"`
			Message string      `json:"message"`
			Error   string      `json:"error"`
		}
		json.NewDecoder(res.Body).Decode(&apiResponse)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, expected.Title, apiResponse.Data.Title)
		assert.Equal(t, expected.Author, apiResponse.Data.Author)
		assert.Equal(t, expected.Year, apiResponse.Data.Year)
		assert.Equal(t, "success", apiResponse.Message)
		assert.Empty(t, apiResponse.Error)
	})

	t.Run("Get non-existing book", func(t *testing.T) {

		app := setupTestApp(t)
		req := httptest.NewRequest("GET", "/books/99", nil)

		res, err := app.Test(req, -1)
		assert.NoError(t, err)

		var apiResponse struct {
			Error   string `json:"error"`
			Message string `json:"message"`
		}
		json.NewDecoder(res.Body).Decode(&apiResponse)

		assert.Equal(t, http.StatusNotFound, res.StatusCode)
		assert.Contains(t, apiResponse.Error, "book not found")
		assert.Equal(t, "error", apiResponse.Message)

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

		var apiResponse struct {
			Error   string `json:"error"`
			Message string `json:"message"`
		}
		json.NewDecoder(res.Body).Decode(&apiResponse)

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		assert.Equal(t, "error", apiResponse.Message)
		assert.Equal(t, "invalid payload", apiResponse.Error)
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

				var apiResponse struct {
					Error   string `json:"error"`
					Message string `json:"message"`
				}
				json.NewDecoder(res.Body).Decode(&apiResponse)

				assert.Equal(t, http.StatusBadRequest, res.StatusCode)
				assert.Equal(t, "error", apiResponse.Message)
				assert.Equal(t, "invalid payload", apiResponse.Error)
			})
		}
	})

	t.Run("creating a book with correct payload", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/books", bytes.NewReader(jsonString))
		req.Header.Add("Content-Type", "application/json")

		res, err := app.Test(req, -1)
		assert.NoError(t, err)

		var apiResponse struct {
			Error   string      `json:"error"`
			Message string      `json:"message"`
			Data    models.Book `json:"data"`
		}
		json.NewDecoder(res.Body).Decode(&apiResponse)

		assert.Equal(t, http.StatusCreated, res.StatusCode)
		assert.Equal(t, book.Title, apiResponse.Data.Title)
		assert.Equal(t, book.Author, apiResponse.Data.Author)
		assert.Equal(t, book.Year, apiResponse.Data.Year)
		assert.Equal(t, "success", apiResponse.Message)
		assert.Empty(t, apiResponse.Error)
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

		var apiResponse struct {
			Error   string `json:"error"`
			Message string `json:"message"`
		}
		json.NewDecoder(res.Body).Decode(&apiResponse)

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		assert.Equal(t, "error", apiResponse.Message)
		assert.Equal(t, "invalid parameter", apiResponse.Error)
	})

	t.Run("updating an existing book without request body", func(t *testing.T) {
		app := setupTestApp(t)

		req := httptest.NewRequest("PUT", "/books/1", nil)
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req, -1)
		assert.NoError(t, err)

		var apiResponse struct {
			Error   string `json:"error"`
			Message string `json:"message"`
		}
		json.NewDecoder(res.Body).Decode(&apiResponse)

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		assert.Equal(t, "error", apiResponse.Message)
		assert.Equal(t, "invalid payload", apiResponse.Error)

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

				var apiResponse struct {
					Error   string `json:"error"`
					Message string `json:"message"`
				}
				json.NewDecoder(res.Body).Decode(&apiResponse)

				assert.Equal(t, http.StatusBadRequest, res.StatusCode)
				assert.Equal(t, "error", apiResponse.Message)
				assert.Equal(t, "invalid payload", apiResponse.Error)
			})
		}

	})

	t.Run("updating an existing book", func(t *testing.T) {
		app := setupTestApp(t)

		req := httptest.NewRequest("PUT", "/books/1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req, -1)
		assert.NoError(t, err)

		var apiResponse struct {
			Error   string      `json:"error"`
			Message string      `json:"message"`
			Data    models.Book `json:"data"`
		}
		json.NewDecoder(res.Body).Decode(&apiResponse)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, uint(1), apiResponse.Data.ID)
		assert.Equal(t, updatedBook.Title, apiResponse.Data.Title)
		assert.Equal(t, updatedBook.Author, apiResponse.Data.Author)
		assert.Equal(t, updatedBook.Year, apiResponse.Data.Year)
		assert.Equal(t, "success", apiResponse.Message)
		assert.Empty(t, apiResponse.Error)
	})

}

func TestDeleteBook(t *testing.T) {
	app := setupTestApp(t)

	t.Run("Deleting a non-existing book", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/books/99", nil)
		res, err := app.Test(req, -1)
		assert.NoError(t, err)

		var apiResponse struct {
			Error   string `json:"error"`
			Message string `json:"message"`
		}
		json.NewDecoder(res.Body).Decode(&apiResponse)

		assert.Equal(t, http.StatusNotFound, res.StatusCode)
		assert.Equal(t, "error", apiResponse.Message)
		assert.Contains(t, apiResponse.Error, "book not found")
	})

	t.Run("Deleting an existing book", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/books/3", nil)
		res, err := app.Test(req, -1)
		assert.NoError(t, err)

		var apiResponse struct {
			Error   string      `json:"error"`
			Message string      `json:"message"`
			Data    models.Book `json:"data"`
		}
		json.NewDecoder(res.Body).Decode(&apiResponse)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "Book Three", apiResponse.Data.Title)
		assert.Equal(t, "Author C", apiResponse.Data.Author)
		assert.Equal(t, 2023, apiResponse.Data.Year)
		assert.Equal(t, "success", apiResponse.Message)
		assert.Empty(t, apiResponse.Error)

		getReq := httptest.NewRequest("GET", "/books/3", nil)
		getRes, err := app.Test(getReq, -1)
		assert.NoError(t, err)

		json.NewDecoder(getRes.Body).Decode(&apiResponse)

		assert.Equal(t, http.StatusNotFound, getRes.StatusCode)
		assert.Equal(t, "error", apiResponse.Message)
		assert.Contains(t, apiResponse.Error, "book not found")

	})

	t.Run("deleting a book with an invalid param", func(t *testing.T) {

		req := httptest.NewRequest("DELETE", "/books/xx", nil)

		res, err := app.Test(req, -1)
		assert.NoError(t, err)

		var apiResponse struct {
			Error   string `json:"error"`
			Message string `json:"message"`
		}
		json.NewDecoder(res.Body).Decode(&apiResponse)

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		assert.Equal(t, "error", apiResponse.Message)
		assert.Equal(t, "invalid parameter", apiResponse.Error)
	})

	t.Run("deleting a book with an empty param", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/books", nil)

		res, err := app.Test(req, -1)
		assert.NoError(t, err)

		var apiResponse struct {
			Error   string `json:"error"`
			Message string `json:"message"`
		}
		json.NewDecoder(res.Body).Decode(&apiResponse)

		assert.Equal(t, http.StatusMethodNotAllowed, res.StatusCode)
		assert.Equal(t, "error", apiResponse.Message)
		assert.Equal(t, "Method Not Allowed", apiResponse.Error)
	})
}

type mockedBookService struct {
	books []*models.Book
}

var _ services.IBookService = (*mockedBookService)(nil)

func NewMockedBookService() *mockedBookService {
	books := []*models.Book{
		{Title: "Book One", Author: "Author A", Year: 2021},
		{Title: "Book Two", Author: "Author B", Year: 2022},
		{Title: "Book Three", Author: "Author C", Year: 2023},
	}

	return &mockedBookService{books: books}
}

func (m *mockedBookService) CreateBook(book *models.Book) (*models.Book, error) {
	book.ID = uint(len(m.books) + 1)
	m.books = append(m.books, book)
	return book, nil
}

func (m *mockedBookService) GetBook(id uint) (*models.Book, error) {
	if id > uint(len(m.books)) {
		return nil, services.ErrNotFound
	}
	return m.books[id-1], nil
}

func (m *mockedBookService) GetAllBooks(page, limit int) ([]*models.Book, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	start := (page - 1) * limit
	end := start + limit

	if start >= len(m.books) {
		return []*models.Book{}, nil
	}

	if end > len(m.books) {
		end = len(m.books)
	}

	return m.books[start:end], nil
}

func (m *mockedBookService) UpdateBook(payload *models.Book) (*models.Book, error) {
	if payload.ID == 0 || payload.ID > uint(len(m.books)) {
		return nil, services.ErrNotFound
	}
	m.books[payload.ID-1] = payload
	return m.books[payload.ID-1], nil
}

func (m *mockedBookService) DeleteBook(id uint) (*models.Book, error) {
	if id == 0 || id > uint(len(m.books)) {
		return nil, services.ErrNotFound
	}
	book := m.books[id-1]
	m.books = append(m.books[:id-1], m.books[id:]...)
	return book, nil
}
