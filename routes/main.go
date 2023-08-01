package routes

import (
	"example/controller"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func InitApp(db *gorm.DB) *echo.Echo {
	e := echo.New()

	bookController := controller.NewBookController(db)
	bookRoute := e.Group("/books")
	bookRoute.POST("", bookController.Create)
	bookRoute.GET("", bookController.Index)
	bookRoute.GET("/:id", bookController.Detail)
	bookRoute.DELETE("/:id", bookController.Delete)

	return e
}
