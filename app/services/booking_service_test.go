package services

import (
	"errors"
	"go-booking-api/app/models"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"syreclabs.com/go/faker"
)

func Test_Book_WhenNoError_Returns_BookingId(t *testing.T) {

	// Arrange
	var db *gorm.DB
	var setupError error

	sqlDb, sqlMock, _ := sqlmock.New()
	defer sqlDb.Close()
	db, setupError = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb,
	}), &gorm.Config{})

	count := uint(faker.RandomInt(0, 100))
	userId := uint(faker.RandomInt(0, 100))
	expectedBookingId := uint(faker.RandomInt(0, 100))

	sqlMock.ExpectQuery(regexp.QuoteMeta("")).
		WithArgs(count, userId).
		WillReturnRows(sqlmock.NewRows([]string{""}).
			AddRow(expectedBookingId))

	assert.Nil(t, setupError)
	target := NewBookingService(db)

	// Act
	actualBookingId, actualError := target.Book(userId, count)

	// Assert
	assert.Nil(t, actualError)
	assert.Equal(t, expectedBookingId, actualBookingId)
}

func Test_Book_WhenDbError_Returns_Error(t *testing.T) {

	// Arrange
	var db *gorm.DB
	var setupError error

	sqlDb, sqlMock, _ := sqlmock.New()
	defer sqlDb.Close()
	db, setupError = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb,
	}), &gorm.Config{})

	count := uint(faker.RandomInt(0, 100))
	userId := uint(faker.RandomInt(0, 100))
	expectedError := errors.New("sb")

	sqlMock.ExpectQuery(regexp.QuoteMeta("")).
		WithArgs(count, userId).
		WillReturnError(expectedError)

	assert.Nil(t, setupError)
	target := NewBookingService(db)

	// Act
	actualBookingId, actualError := target.Book(userId, count)

	// Assert
	assert.Empty(t, actualBookingId)
	assert.Equal(t, expectedError, actualError)
}

func Test_Book_WhenDbBookingIdZero_Returns_Error(t *testing.T) {

	// Arrange
	var db *gorm.DB
	var setupError error

	sqlDb, sqlMock, _ := sqlmock.New()
	defer sqlDb.Close()
	db, setupError = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb,
	}), &gorm.Config{})

	count := uint(faker.RandomInt(0, 100))
	userId := uint(faker.RandomInt(0, 100))
	expectedError := errors.New("sb")

	sqlMock.ExpectQuery(regexp.QuoteMeta("")).
		WithArgs(count, userId).
		WillReturnRows(sqlmock.NewRows([]string{""}).
			AddRow(0))

	assert.Nil(t, setupError)
	target := NewBookingService(db)

	// Act
	actualBookingId, actualError := target.Book(userId, count)

	// Assert
	assert.Empty(t, actualBookingId)
	assert.Equal(t, expectedError.Error(), actualError.Error())
}

func Test_GetBookings_WhenNoError_Returns_Bookings(t *testing.T) {

	// Arrange
	var db *gorm.DB
	var setupError error

	sqlDb, sqlMock, _ := sqlmock.New()
	defer sqlDb.Close()
	db, setupError = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb,
	}), &gorm.Config{})

	userId := uint(faker.RandomInt(0, 100))
	expectedBookings := []models.Booking{{UserID: 1, TicketInventoryID: 1, Tickets: 16}}

	sqlMock.ExpectQuery(regexp.QuoteMeta("")).
		WithArgs(userId).
		WillReturnRows(sqlmock.NewRows([]string{"UserID", "TicketInventoryID", "Tickets"}).
			AddRow(expectedBookings[0].UserID, expectedBookings[0].TicketInventoryID, expectedBookings[0].Tickets))

	assert.Nil(t, setupError)
	target := NewBookingService(db)

	// Act
	actualBookings, actualError := target.GetBookings(userId)

	// Assert
	assert.Equal(t, &expectedBookings, actualBookings)
	assert.Nil(t, actualError)
}

func Test_GetBookings_WhenDbError_Returns_Error(t *testing.T) {

	// Arrange
	var db *gorm.DB
	var setupError error

	sqlDb, sqlMock, _ := sqlmock.New()
	defer sqlDb.Close()
	db, setupError = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb,
	}), &gorm.Config{})

	userId := uint(faker.RandomInt(0, 100))
	expectedError := errors.New("sb")

	sqlMock.ExpectQuery(regexp.QuoteMeta("")).
		WithArgs(userId).
		WillReturnError(expectedError)

	assert.Nil(t, setupError)
	target := NewBookingService(db)

	// Act
	actualBookings, actualError := target.GetBookings(userId)

	// Assert
	assert.Empty(t, actualBookings)
	assert.Equal(t, expectedError, actualError)
}
