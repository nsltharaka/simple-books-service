package services

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/nsltharaka/booksapi/models"
	"gorm.io/gorm"
)

var (
	ErrNotFound = errors.New("book not found")
)

type BookService struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewBookService(db *gorm.DB, logger *slog.Logger) *BookService {
	return &BookService{db: db, logger: logger}
}

func (s *BookService) CreateBook(book *models.Book) (*models.Book, error) {
	if err := s.db.Create(book).Error; err != nil {
		s.logger.Error("failed to create new book", "error", err)
		return nil, fmt.Errorf("failed to create new book : %w", err)
	}
	s.logger.Info("created new book", "book", book)
	return book, nil
}

func (s *BookService) GetBook(id uint) (*models.Book, error) {
	var book models.Book
	if err := s.db.First(&book, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("book not found", "id", id)
			return nil, fmt.Errorf("%w: id %d", ErrNotFound, id)
		}
		s.logger.Error("error fetching book", "id", id, "error", err)
		return nil, fmt.Errorf("error while fetching the book : %w", err)
	}
	s.logger.Info("fetched book", "book", book)
	return &book, nil
}

func (s *BookService) GetAllBooks(page, limit int) ([]*models.Book, error) {
	var books []*models.Book
	offset := (page - 1) * limit
	if err := s.db.Limit(limit).Offset(offset).Find(&books).Error; err != nil {
		s.logger.Error("error fetching paginated books", "error", err)
		return nil, fmt.Errorf("error while fetching books : %w", err)
	}
	s.logger.Info("fetched paginated books", "count", len(books), "page", page, "limit", limit)
	return books, nil
}

func (s *BookService) UpdateBook(payload *models.Book) (*models.Book, error) {
	var book models.Book
	if err := s.db.First(&book, payload.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("book to update not found", "id", payload.ID)
			return nil, fmt.Errorf("%w: id %d", ErrNotFound, payload.ID)
		}
		s.logger.Error("error fetching book for update", "id", payload.ID, "error", err)
		return nil, fmt.Errorf("error while fetching the book : %w", err)
	}

	book.Title = payload.Title
	book.Author = payload.Author
	book.Year = payload.Year

	if err := s.db.Save(&book).Error; err != nil {
		s.logger.Error("error saving updated book", "book", book, "error", err)
		return nil, fmt.Errorf("error while saving the book : %w", err)
	}
	s.logger.Info("updated book", "book", book)
	return &book, nil
}

func (s *BookService) DeleteBook(id uint) (*models.Book, error) {
	var book models.Book
	if err := s.db.First(&book, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("book to delete not found", "id", id)
			return nil, fmt.Errorf("%w: id %d", ErrNotFound, id)
		}
		s.logger.Error("error fetching book for deletion", "id", id, "error", err)
		return nil, fmt.Errorf("error while fetching the book : %w", err)
	}

	if err := s.db.Delete(&book).Error; err != nil {
		s.logger.Error("error deleting book", "book", book, "error", err)
		return nil, fmt.Errorf("error while deleting the book : %w", err)
	}
	s.logger.Info("deleted book", "book", book)
	return &book, nil
}
