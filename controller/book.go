package controller

import (
	"example/model"
	"example/repository"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type BookController struct {
	Repository repository.BookRepository
}

func NewBookController(r repository.BookRepository) BookController {
	return BookController{Repository: r}
}

func (c BookController) Create(ctx echo.Context) error {
	newBook := model.Book{}
	if err := ctx.Bind(&newBook); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	if err := c.Repository.Create(&newBook); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	response := map[string]interface{}{
		"message": "success create book",
		"newBook": newBook,
	}

	return ctx.JSON(http.StatusCreated, response)
}

func (c BookController) Detail(ctx echo.Context) error {
	id := ctx.Param("id")
	idInt, _ := strconv.Atoi(id)
	book, err := c.Repository.FindById(idInt)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"message": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, book)
}

func (c BookController) Index(ctx echo.Context) error {
	books, _ := c.Repository.FindAll()

	return ctx.JSON(http.StatusOK, books)
}

func (c BookController) Delete(ctx echo.Context) error {
	id := ctx.Param("id")
	idInt, _ := strconv.Atoi(id)

	if err := c.Repository.Delete(idInt); err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"message": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "success delete book",
	})
}
