package services

import (
	"fmt"
	"go-booking-api/app/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Test_InitDb_CreatesDb_If_NotExist(t *testing.T) {

	// Arrange
	const DB_NAME = "go-booking-api-init-db-test"
	createDb := fmt.Sprintf(SQL_STATEMENT_CREATE_DB, DB_NAME)
	dbInfo := DbInfo{Host: "localhost", Port: 5432, User: "postgres", Password: "postgres", DbName: DB_NAME, SslMode: "disable", TimeZone: "Europe/London"}
	var db *gorm.DB
	target := NewDBService(db)

	// Act
	db = target.InitDB(dbInfo, gorm.Open, createDb)

	// Assert
	assert.Nil(t, db.Error)
	assert.True(t, db.Migrator().HasTable(&models.TicketInventory{}))
	assert.True(t, db.Migrator().HasTable(&models.Booking{}))
	assert.True(t, db.Migrator().HasTable(&models.Credentials{}))
	assert.True(t, db.Migrator().HasTable(&models.User{}))

	sqlDb, _ := db.DB()
	sqlDb.Close()

	cleanUp(dbInfo)
}

func cleanUp(dbInfo DbInfo) {

	conn := fmt.Sprintf("host=%s user=%s password=%s port=%d sslmode=%s TimeZone=%s", dbInfo.Host, dbInfo.User, dbInfo.Password, dbInfo.Port, dbInfo.SslMode, dbInfo.TimeZone)

	db2, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  conn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	if err != nil {

		panic(fmt.Sprintf("Manually delete %s database", dbInfo.DbName))
	}

	dropDb := fmt.Sprintf(`drop database "%s"`, dbInfo.DbName)
	db2 = db2.Exec(dropDb)

	if db2.Error != nil {

		panic(fmt.Sprintf("Manually delete %s database", dbInfo.DbName))
	}
}
