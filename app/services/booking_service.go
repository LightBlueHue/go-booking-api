package services

type IBookingService interface {
	Book(count uint)
}

type BookingService struct {
}

func GetBookingService() IBookingService {
	return &BookingService{}
}

func (s *BookingService) Book(count uint) {
	// db := GetDBService().GetDB()
	
	// var pwd string
	//db.Raw("SELECT password FROM users INNER JOIN credentials ON users.credential_id = credentials.id WHERE email = ?", email).Scan(&pwd)
}
