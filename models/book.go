package models

type Book struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	Title  string `json:"title" validate:"required,endsnotwith= "`
	Author string `json:"author" validate:"required,endsnotwith= "`
	Year   int    `json:"year" validate:"required,number"`
}
