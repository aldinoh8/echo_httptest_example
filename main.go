package main

import (
	"example/config"
	"example/controller"

	"github.com/labstack/echo/v4"
)

func main() {
	dsn := "host=localhost user=postgres password=postgres dbname=sandbox_db port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	e := echo.New()

	// Connect To Database
	db := config.InitDatabase(dsn)
	bookController := controller.NewBookController(db)

	bookRoute := e.Group("/books")
	bookRoute.POST("", bookController.Create)
	bookRoute.GET("", bookController.Index)
	bookRoute.GET("/:id", bookController.Detail)
	bookRoute.DELETE("/:id", bookController.Delete)

	e.Logger.Fatal(e.Start(":8000"))
}
