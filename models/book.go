package models

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	Title  string `json:"title" validate:"required,endsnotwith= "`
	Author string `json:"author" validate:"required,endsnotwith= "`
	Year   int    `json:"year" validate:"required,number"`
}
