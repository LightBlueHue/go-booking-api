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
	MATCH_ONE_OR_MORE_NON_WHITESPACE_CHARACTER    = "^\\S+$"
	MATCH_ANY_WORD_CHARACTER_OF_LENGTH_4_TO_20    = "^\\w{4,20}$"
)

type ValidationService struct {
}

func NewValidationService() IValidationService {

	return &ValidationService{}
}

// ValidateLoginRequest validates login request model.
func (s *ValidationService) ValidateLoginRequest(v *revel.Validation, l *requests.LoginRequest) {

	v.Email(l.Email)
	v.Match(l.Password, regexp.MustCompile(MATCH_ANY_WORD_CHARACTER_OF_LENGTH_4_TO_20)).Message(VALIDATION_REQUEST_PASSWORD_LENGTH)
}

// ValidateRegisterRequest validates register request mode.
func (s *ValidationService) ValidateRegisterRequest(v *revel.Validation, l *requests.RegisterRequest) {

	v.Match(l.FirstName, regexp.MustCompile(MATCH_ONE_OR_MORE_NON_WHITESPACE_CHARACTER)).Message(VALIDATION_REGISTER_REQUEST_FIRSTNAME_INVALID)
	v.Match(l.LastName, regexp.MustCompile(MATCH_ONE_OR_MORE_NON_WHITESPACE_CHARACTER)).Message(VALIDATION_REGISTER_REQUEST_LASTNAME_INVALID)
	v.Email(l.Email)
	v.Match(l.Password, regexp.MustCompile(MATCH_ANY_WORD_CHARACTER_OF_LENGTH_4_TO_20)).Message(VALIDATION_REQUEST_PASSWORD_LENGTH)
	v.Required(l.Password == l.ConfirmPassword).Message(VALIDATION_REGISTER_REQUEST_PASSWORD_NOMATCH)
}

// ValidateBookingRequest validates book request model
func (s *ValidationService) ValidateBookingRequest(v *revel.Validation, count uint) {

	min := int(count)
	v.Min(min, 1).Message(VALIDATION_BOOKING_REQUEST_VALID_NUMBER)
}
