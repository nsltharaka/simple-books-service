package services

import (
	"errors"
	"fmt"

	"github.com/nsltharaka/booksapi/models"
	"gorm.io/gorm"
)

var (
	ErrNotFound = errors.New("book not found")
)

type BookService struct {
	db *gorm.DB
}

func NewBookService(db *gorm.DB) *BookService {
	return &BookService{db: db}
}

func (s *BookService) CreateBook(book *models.Book) (*models.Book, error) {
	if err := s.db.Create(book).Error; err != nil {
		return nil, fmt.Errorf("failed to create new book : %w", err)
	}

	return book, nil
}

func (s *BookService) GetBook(id uint) (*models.Book, error) {
	var book models.Book
	if err := s.db.First(&book, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: id %d", ErrNotFound, id)
		}
		return nil, fmt.Errorf("error while fetching the book : %w", err)
	}

	return &book, nil
}

func (s *BookService) GetAllBooks() ([]models.Book, error) {
	var books []models.Book
	if err := s.db.Find(&books).Error; err != nil {
		return nil, fmt.Errorf("error while fetching books : %w", err)
	}

	return books, nil
}

func (s *BookService) UpdateBook(payload *models.Book) (*models.Book, error) {
	var book models.Book
	if err := s.db.First(&book, payload.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: id %d", ErrNotFound, payload.ID)
		}
		return nil, fmt.Errorf("error while fetching the book : %w", err)
	}

	book.Title = payload.Title
	book.Author = payload.Author
	book.Year = payload.Year

	if err := s.db.Save(&book).Error; err != nil {
		return nil, fmt.Errorf("error while saving the book : %w", err)
	}
	return &book, nil
}

func (s *BookService) DeleteBook(id uint) (*models.Book, error) {
	var book models.Book
	if err := s.db.First(&book, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: id %d", ErrNotFound, id)
		}
		return nil, fmt.Errorf("error while fetching the book : %w", err)
	}

	if err := s.db.Delete(&book).Error; err != nil {
		return nil, fmt.Errorf("error while deleting the book : %w", err)
	}
	return &book, nil
}
