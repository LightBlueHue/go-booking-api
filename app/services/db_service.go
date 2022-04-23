package services

import (
	"go-booking-api/app/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type IDBService interface{
	InitDB()
	GetDB() *gorm.DB
}

type DBService struct {
}

var db *gorm.DB

func GetDBService() IDBService {

	return &DBService{}
}

func (dbService *DBService) InitDB() {

	database, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "user=postgres password=postgres dbname=go-booking-api port=5432 sslmode=disable TimeZone=Europe/London",
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	database.AutoMigrate(&models.Booking{})
	database.AutoMigrate(&models.Credentials{})
	database.AutoMigrate(&models.User{})

	db = database
}

func (dbService *DBService) GetDB() *gorm.DB {

	if db == nil {
		dbService.InitDB()
	}
	return db
}
