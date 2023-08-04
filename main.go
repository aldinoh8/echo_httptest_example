package main

import (
	"example/config"
	"example/routes"
)

func main() {
	dsn := "host=localhost user=postgres password=postgres dbname=sandbox_db port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	db := config.InitDatabase(dsn)
	e := routes.InitApp(db)

	e.Logger.Fatal(e.Start(":8000"))
}
