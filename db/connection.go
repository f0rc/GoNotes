package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"samplegoapp.com/models"
)

var DB *gorm.DB

func Connect() {
	dsn := "host=localhost user=postgres password=8TQDx2%qEk6]]{sK dbname=blogapp port=5432 sslmode=disable"
	connection, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}
	DB = connection

	connection.AutoMigrate(&models.User{})
}
