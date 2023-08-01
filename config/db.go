package config

import (
	"example/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDatabase(dbDsn string) *gorm.DB {
	database, err := gorm.Open(postgres.Open(dbDsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	database.AutoMigrate(&model.Book{})
	return database
}
