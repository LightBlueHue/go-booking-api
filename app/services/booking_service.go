package services

import (
	"fmt"
	"go-booking-api/app/models"
)

type IBookingService interface {
	Book(user *models.User, count uint) (uint, error)
}

type BookingService struct {
}

func GetBookingService() IBookingService {

	return &BookingService{}
}

func (s *BookingService) Book(user *models.User, count uint) (uint, error) {

	db := GetDBService().GetDB()

	var bookingId uint
	result := db.Raw("select book(?,?);", count, user.ID).Scan(&bookingId)

	if bookingId == 0 || result.RowsAffected == 0 {

		return bookingId, fmt.Errorf("sb")
	}

	return bookingId, result.Error
}
