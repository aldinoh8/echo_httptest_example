package repository

import (
	"errors"
	"example/model"

	"gorm.io/gorm"
)

type BookRepository interface {
	FindAll() ([]model.Book, error)
	FindById(int) (model.Book, error)
	Create(*model.Book) error
	Delete(int) error
}

type Book struct {
	DB *gorm.DB
}

func NewBook(db *gorm.DB) Book {
	return Book{DB: db}
}

func (br Book) FindAll() (books []model.Book, err error) {
	result := br.DB.Find(&books)
	if result.Error != nil {
		panic(err)
	}

	return books, err
}

func (br Book) FindById(id int) (book model.Book, err error) {
	result := br.DB.First(&book, id)
	if result.RowsAffected == 0 {
		return book, errors.New("book not found")
	}
	return book, err
}

func (br Book) Create(book *model.Book) (err error) {
	result := br.DB.Create(book)
	if result.Error != nil {
		return errors.New("failed to create new book")
	}

	return nil
}

func (br Book) Delete(id int) (err error) {
	result := br.DB.Delete(&model.Book{}, id)
	if result.RowsAffected == 0 {
		return errors.New("book not found")
	}

	if result.Error != nil {
		panic(result.Error)
	}

	return err
}
