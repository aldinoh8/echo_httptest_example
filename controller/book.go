package controller

import (
	"example/model"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type BookController struct {
	DB *gorm.DB
}

func NewBookController(db *gorm.DB) BookController {
	return BookController{DB: db}
}

func (c BookController) Create(ctx echo.Context) error {
	newBook := model.Book{}
	if err := ctx.Bind(&newBook); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	if err := c.DB.Create(&newBook).Error; err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	response := map[string]interface{}{
		"message": "success create book",
		"newBook": newBook,
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c BookController) Detail(ctx echo.Context) error {
	id := ctx.Param("id")
	book := model.Book{}

	result := c.DB.First(&book, id)
	if result.RowsAffected == 0 {
		return ctx.JSON(http.StatusNotFound, map[string]string{"Message": "Book Not Found"})
	}
	if result.Error != nil {
		panic(result.Error)
	}

	return ctx.JSON(http.StatusOK, book)
}

func (c BookController) Index(ctx echo.Context) error {
	books := []model.Book{}
	if err := c.DB.Find(&books).Error; err != nil {
		panic(err)
	}

	return ctx.JSON(http.StatusOK, books)
}

func (c BookController) Delete(ctx echo.Context) error {
	id := ctx.Param("id")

	result := c.DB.Delete(&model.Book{}, id)
	if result.RowsAffected == 0 {
		return ctx.JSON(http.StatusNotFound, map[string]string{"Message": "Book Not Found"})
	}
	if result.Error != nil {
		panic(result.Error)
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"message": "Success delete book",
	})
}
