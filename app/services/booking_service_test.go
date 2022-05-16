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
	target := GetBookingService(db)

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
	target := GetBookingService(db)

	actualBookingId, actualError := target.Book(userId, count)

	assert.Empty(t, actualBookingId)
	assert.Equal(t, expectedError, actualError)
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
	target := GetBookingService(db)

	actualBookingId, actualError := target.Book(userId, count)

	assert.Empty(t, actualBookingId)
	assert.Equal(t, expectedError, actualError)
}

func Test_GetBookings_Returns_Data(t *testing.T) {

	var db *gorm.DB
	var setupError error

	sqlDb, sqlMock, _ := sqlmock.New()
	defer sqlDb.Close()
	userId := uint(faker.RandomInt(0, 100))
	expectedBookings := []models.Booking{{UserID: 1, TicketInventoryID: 1, Tickets: 16}}

	sqlMock.ExpectQuery(regexp.QuoteMeta("user_id = $1")).
		WithArgs(userId).
		WillReturnRows(sqlmock.NewRows([]string{"UserID", "TicketInventoryID", "Tickets"}).
			AddRow(expectedBookings[0].UserID, expectedBookings[0].TicketInventoryID, expectedBookings[0].Tickets))

	db, setupError = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb,
	}), &gorm.Config{})

	assert.Nil(t, setupError)
	target := GetBookingService(db)

	actualBookings, actualError := target.GetBookings(userId)

	assert.Equal(t, &expectedBookings, actualBookings)
	assert.Nil(t, actualError)
}

func Test_GetBookings_Returns_Error(t *testing.T) {

	var db *gorm.DB
	var setupError error

	sqlDb, sqlMock, _ := sqlmock.New()
	defer sqlDb.Close()
	userId := uint(faker.RandomInt(0, 100))
	expectedError := errors.New("sb")

	sqlMock.ExpectQuery(regexp.QuoteMeta("user_id = $1")).
		WithArgs(userId).
		WillReturnError(expectedError)

	db, setupError = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb,
	}), &gorm.Config{})

	assert.Nil(t, setupError)
	target := GetBookingService(db)

	actualBookings, actualError := target.GetBookings(userId)

	assert.Empty(t, actualBookings)
	assert.Equal(t, expectedError, actualError)
}
