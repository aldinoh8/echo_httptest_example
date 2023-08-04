package routes

import (
	"example/controller"
	"example/repository"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func InitApp(db *gorm.DB) *echo.Echo {
	e := echo.New()

	bookRepository := repository.NewBook(db)
	bookController := controller.NewBookController(bookRepository)

	bookRoute := e.Group("/books")
	bookRoute.POST("", bookController.Create)
	bookRoute.GET("", bookController.Index)
	bookRoute.GET("/:id", bookController.Detail)
	bookRoute.DELETE("/:id", bookController.Delete)

	return e
}
