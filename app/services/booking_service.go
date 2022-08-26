// Package services contains helper classes.
package services

import (
	"fmt"
	"go-booking-api/app/models"

	"gorm.io/gorm"
)

const (
	SQL_STATEMENT_CALL_BOOK_FUNCTION      = "select book(?,?);"
	SQL_STATEMENT_GET_BOOKINGS_BY_USER_ID = "user_id = ?"
)

type BookingService struct {
	db *gorm.DB
}

func NewBookingService(db *gorm.DB) IBookingService {

	return &BookingService{db}
}

// Book books a ticket for the associated user id.
func (s *BookingService) Book(userId uint, count uint) (uint, error) {

	var bookingId uint
	result := s.db.Raw(SQL_STATEMENT_CALL_BOOK_FUNCTION, count, userId).Scan(&bookingId)

	if bookingId == 0 || result.RowsAffected == 0 {

		return bookingId, fmt.Errorf("sb")
	}

	return bookingId, result.Error
}

// GetBookings returns tickets booked for the associated user id.
func (s *BookingService) GetBookings(userId uint) (*[]models.Booking, error) {

	var bookings *[]models.Booking
	result := s.db.Where(SQL_STATEMENT_GET_BOOKINGS_BY_USER_ID, userId).Find(&bookings)
	return bookings, result.Error
}
