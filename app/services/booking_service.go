package services

import (
	"fmt"
	"go-booking-api/app/models"
)

const (
	SQL_STATEMENT_CALL_BOOK_FUNCTION      = "select book(?,?);"
	SQL_STATEMENT_GET_BOOKINGS_BY_USER_ID = "user_id = ?"
)

type IBookingService interface {
	Book(user *models.User, count uint) (uint, error)
	GetBookings(user *models.User) (*[]models.Booking, error)
}

type BookingService struct {
}

func GetBookingService() IBookingService {

	return &BookingService{}
}

func (s *BookingService) Book(user *models.User, count uint) (uint, error) {

	db := GetDBService().GetDB()

	var bookingId uint
	result := db.Raw(SQL_STATEMENT_CALL_BOOK_FUNCTION, count, user.ID).Scan(&bookingId)

	if bookingId == 0 || result.RowsAffected == 0 {

		return bookingId, fmt.Errorf("sb")
	}

	return bookingId, result.Error
}

func (s *BookingService) GetBookings(user *models.User) (*[]models.Booking, error) {

	db := GetDBService().GetDB()
	var bookings *[]models.Booking
	result := db.Where(SQL_STATEMENT_GET_BOOKINGS_BY_USER_ID, user.ID).Find(&bookings)
	return bookings, result.Error
}
