package services

import (
	"log/slog"
	"os"
	"testing"

	"github.com/nsltharaka/booksapi/database"
	"github.com/nsltharaka/booksapi/models"
	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) (*BookService, func()) {
	filename := "test_service.db"
	os.Setenv("SQLITE_FILENAME", filename)
	db, err := database.Connect()
	if err != nil {
		t.Fatalf("failed to connect to test db: %v", err)
	}

	books := []models.Book{
		{Title: "Book One", Author: "Author A", Year: 2021},
		{Title: "Book Two", Author: "Author B", Year: 2022},
		{Title: "Book Three", Author: "Author C", Year: 2023},
	}

	service := NewBookService(db, slog.New(slog.NewTextHandler(os.Stdout, nil)))

	for _, book := range books {
		service.CreateBook(&models.Book{
			Title:  book.Title,
			Author: book.Author,
			Year:   book.Year,
		})
	}

	cleanup := func() {
		os.Unsetenv("SQLITE_FILENAME")
		os.Remove(filename)
	}

	return service, cleanup
}

func TestCreateAndGetBook(t *testing.T) {
	service, cleanup := setupTestDB(t)
	t.Cleanup(cleanup)
	book := &models.Book{Title: "Test Book", Author: "Author", Year: 2023}

	createdBook, err := service.CreateBook(book)
	assert.NoError(t, err)
	assert.NotZero(t, createdBook.ID)

	fetchedBook, err := service.GetBook(createdBook.ID)
	assert.NoError(t, err)
	assert.Equal(t, uint(4), fetchedBook.ID)
	assert.Equal(t, createdBook.ID, fetchedBook.ID)
	assert.Equal(t, createdBook.Title, fetchedBook.Title)
}

func TestGetBooks(t *testing.T) {

	service, cleanUp := setupTestDB(t)
	t.Cleanup(cleanUp)

	t.Run("fetching an existing book", func(t *testing.T) {
		bookId := 2
		want := models.Book{Title: "Book Two", Author: "Author B", Year: 2022}
		fetchedBook, err := service.GetBook(uint(bookId))
		assert.NoError(t, err)
		assert.Equal(t, uint(bookId), fetchedBook.ID)
		assert.Equal(t, want.Title, fetchedBook.Title)
		assert.Equal(t, want.Author, fetchedBook.Author)
		assert.Equal(t, want.Year, fetchedBook.Year)
	})

	t.Run("fetching non-existing book", func(t *testing.T) {
		var want uint = 99
		fetchedBook, err := service.GetBook(want)
		assert.Error(t, err)
		assert.Nil(t, fetchedBook)
	})

	t.Run("fetching all books", func(t *testing.T) {
		fetchedBooks, err := service.GetAllBooks(1, 10)
		assert.NoError(t, err)
		assert.Equal(t, 3, len(fetchedBooks))
	})

}

func TestUpdateBook(t *testing.T) {
	service, cleanup := setupTestDB(t)
	defer cleanup()

	book := &models.Book{Title: "Old Title", Author: "Author", Year: 2022}
	createdBook, _ := service.CreateBook(book)

	t.Run("updating an existing book", func(t *testing.T) {
		createdBook.Title = "New Title"
		updatedBook, err := service.UpdateBook(createdBook)
		assert.NoError(t, err)
		assert.Equal(t, "New Title", updatedBook.Title)
	})

	t.Run("updating non-existing book", func(t *testing.T) {
		createdBook.ID = 99
		updatedBook, err := service.UpdateBook(createdBook)
		assert.Nil(t, updatedBook)
		assert.Error(t, err)
	})

}

func TestDeleteBook(t *testing.T) {
	service, cleanup := setupTestDB(t)
	defer cleanup()

	book := &models.Book{Title: "To be deleted", Author: "Author", Year: 2021}
	createdBook, _ := service.CreateBook(book)

	t.Run("deleting an existing book", func(t *testing.T) {
		deletedBook, err := service.DeleteBook(createdBook.ID)
		assert.NoError(t, err)
		assert.Equal(t, createdBook.ID, deletedBook.ID)

		_, err = service.GetBook(createdBook.ID)
		assert.Error(t, err)
	})

	t.Run("updating non-existing book", func(t *testing.T) {
		createdBook.ID = 99
		deletedBook, err := service.DeleteBook(createdBook.ID)
		assert.Nil(t, deletedBook)
		assert.Error(t, err)
	})

}
