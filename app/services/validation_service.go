package services

import (
	"go-booking-api/app/models/requests"
	"regexp"

	"github.com/revel/revel"
)

type IValidationService interface {
	ValidateLoginRequest(c *revel.Controller, l *requests.LoginRequest)
	ValidateRegisterRequest(c *revel.Controller, l *requests.RegisterRequest)
}

type ValidationService struct {
}

func GetValidationService() IValidationService {

	return &ValidationService{}
}

func (s *ValidationService) ValidateLoginRequest(c *revel.Controller, l *requests.LoginRequest) {

	c.Validation.Email(l.Email)
	c.Validation.Match(l.Password, regexp.MustCompile("^\\w{4,20}$")).Message("Password must be 4 to 20 characters long")
}

func (s *ValidationService) ValidateRegisterRequest(c *revel.Controller, l *requests.RegisterRequest) {

	c.Validation.Match(l.FirstName, regexp.MustCompile("^\\S+$")).Message("First name invalid")
	c.Validation.Match(l.LastName, regexp.MustCompile("^\\S+$")).Message("Last name invalid")
	c.Validation.Email(l.Email)
	c.Validation.Match(l.Password, regexp.MustCompile("^\\w{4,20}$")).Message("Password must be 4 to 20 characters long")
	c.Validation.Required(l.Password == l.ConfirmPassword).Message("Passwords do not match")
}
