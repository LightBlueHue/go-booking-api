package tests

import (
	"errors"
	"go-booking-api/app/services"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"syreclabs.com/go/faker"
)

func Test_Book_Returns_BookingId(t *testing.T) {

	var db *gorm.DB
	var setupError error

	sqlDb, sqlMock, _ := sqlmock.New()
	defer sqlDb.Close()
	count := uint(faker.RandomInt(0, 100))
	userId := uint(faker.RandomInt(0, 100))
	expectedBookingId := uint(faker.RandomInt(0, 100))

	sqlMock.ExpectQuery(regexp.QuoteMeta("select book($1,$2);")).
		WithArgs(count, userId).
		WillReturnRows(sqlmock.NewRows([]string{"bookingId"}).
			AddRow(expectedBookingId))

	db, setupError = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb,
	}), &gorm.Config{})

	assert.Nil(t, setupError)
	target := services.GetBookingService(db)

	actualBookingId, actualError := target.Book(userId, count)

	assert.Nil(t, actualError)
	assert.Equal(t, expectedBookingId, actualBookingId)
}

func Test_Book_WhenDbError_Returns_Error(t *testing.T) {

	var db *gorm.DB
	var setupError error

	sqlDb, sqlMock, _ := sqlmock.New()
	defer sqlDb.Close()
	count := uint(faker.RandomInt(0, 100))
	userId := uint(faker.RandomInt(0, 100))
	expectedError := errors.New("sb")

	sqlMock.ExpectQuery(regexp.QuoteMeta("select book($1,$2);")).
		WithArgs(count, userId).
		WillReturnError(expectedError)

	db, setupError = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb,
	}), &gorm.Config{})

	assert.Nil(t, setupError)
	target := services.GetBookingService(db)

	actualBookingId, actualError := target.Book(userId, count)

	assert.Empty(t, actualBookingId)
	assert.Equal(t, expectedError.Error(), actualError.Error())
}

func Test_Book_WhenDbBookingIdZero_Returns_Error(t *testing.T) {

	var db *gorm.DB
	var setupError error

	sqlDb, sqlMock, _ := sqlmock.New()
	defer sqlDb.Close()
	count := uint(faker.RandomInt(0, 100))
	userId := uint(faker.RandomInt(0, 100))
	expectedError := errors.New("sb")

	sqlMock.ExpectQuery(regexp.QuoteMeta("select book($1,$2);")).
		WithArgs(count, userId).
		WillReturnRows(sqlmock.NewRows([]string{"bookingId"}).
			AddRow(0))

	db, setupError = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb,
	}), &gorm.Config{})

	assert.Nil(t, setupError)
	target := services.GetBookingService(db)

	actualBookingId, actualError := target.Book(userId, count)

	assert.Empty(t, actualBookingId)
	assert.Equal(t, expectedError.Error(), actualError.Error())
}
