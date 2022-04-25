package services

import "go-booking-api/app/models"

type IBookingService interface {
	Book(user *models.User, count uint) error
}

type BookingService struct {
}

func GetBookingService() IBookingService {
	return &BookingService{}
}

func (s *BookingService) Book(user *models.User, count uint) error {
	// db := GetDBService().GetDB()

	// var pwd string
	//db.Raw("SELECT password FROM users INNER JOIN credentials ON users.credential_id = credentials.id WHERE email = ?", email).Scan(&pwd)

	return nil
}
