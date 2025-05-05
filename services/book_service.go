package services

import (
	"github.com/nsltharaka/booksapi/models"
	"gorm.io/gorm"
)

type BookService struct {
	db *gorm.DB
}

func NewBookService(db *gorm.DB) *BookService {
	return &BookService{db: db}
}

func (s *BookService) CreateBook(book *models.Book) (*models.Book, error) {
	if err := s.db.Create(book).Error; err != nil {
		return nil, err
	}

	return book, nil
}

func (s *BookService) GetBook(id uint) (*models.Book, error) {
	var book models.Book
	if err := s.db.First(&book, id).Error; err != nil {
		return nil, err
	}

	return &book, nil
}

func (s *BookService) GetAllBooks() ([]models.Book, error) {
	var books []models.Book
	if err := s.db.Find(&books).Error; err != nil {
		return nil, err
	}

	return books, nil
}

func (s *BookService) UpdateBook(payload *models.Book) (*models.Book, error) {
	var book models.Book
	if err := s.db.First(&book, payload.ID).Error; err != nil {
		return nil, err
	}

	book.Title = payload.Title
	book.Author = payload.Author
	book.Year = payload.Year

	if err := s.db.Save(&book).Error; err != nil {
		return nil, err
	}
	return &book, nil
}

func (s *BookService) DeleteBook(id uint) (*models.Book, error) {
	var book models.Book
	if err := s.db.First(&book, id).Error; err != nil {
		return nil, err
	}

	if err := s.db.Delete(&book).Error; err != nil {
		return nil, err
	}
	return &book, nil
}
