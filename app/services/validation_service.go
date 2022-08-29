package services

import (
	"go-booking-api/app/models/requests"
	"regexp"

	"github.com/revel/revel"
)

const (
	VALIDATION_REQUEST_PASSWORD_LENGTH            = "Password must be 4 to 20 characters long"
	VALIDATION_REGISTER_REQUEST_FIRSTNAME_INVALID = "First name invalid"
	VALIDATION_REGISTER_REQUEST_LASTNAME_INVALID  = "Last name invalid"
	VALIDATION_REGISTER_REQUEST_PASSWORD_NOMATCH  = "Passwords do not match"
	VALIDATION_BOOKING_REQUEST_VALID_NUMBER       = "Booking must be a positive number, 1 or more"
)

type ValidationService struct {
	*revel.Validation
}

func NewValidationService(v *revel.Validation) IValidationService {

	return &ValidationService{v}
}

func (s *ValidationService) ValidateLoginRequest(l *requests.LoginRequest) {

	s.Email(l.Email)
	s.Match(l.Password, regexp.MustCompile("^\\w{4,20}$")).Message(VALIDATION_REQUEST_PASSWORD_LENGTH)
}

func (s *ValidationService) ValidateRegisterRequest(l *requests.RegisterRequest) {

	s.Match(l.FirstName, regexp.MustCompile("^\\S+$")).Message(VALIDATION_REGISTER_REQUEST_FIRSTNAME_INVALID)
	s.Match(l.LastName, regexp.MustCompile("^\\S+$")).Message(VALIDATION_REGISTER_REQUEST_LASTNAME_INVALID)
	s.Email(l.Email)
	s.Match(l.Password, regexp.MustCompile("^\\w{4,20}$")).Message(VALIDATION_REQUEST_PASSWORD_LENGTH)
	s.Required(l.Password == l.ConfirmPassword).Message(VALIDATION_REGISTER_REQUEST_PASSWORD_NOMATCH)
}

func (v *ValidationService) ValidateBookingRequest(count uint) {

	min := int(count)
	v.Min(min, 1).Message(VALIDATION_BOOKING_REQUEST_VALID_NUMBER)
}
