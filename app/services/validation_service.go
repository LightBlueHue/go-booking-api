package services

import (
	"go-booking-api/app/models/requests"

	"github.com/revel/revel"
)

func ValidateLoginRequest(c *revel.Controller, l *requests.LoginRequest) {

	c.Validation.Required(l.Email)
	c.Validation.Email(l.Email)
	c.Validation.MinSize(l.Password, 4).Message("Password must more than 4 characters long")
}

func ValidateRegisterRequest(c *revel.Controller, l *requests.RegisterRequest) {

	c.Validation.Required(l.Email)
	c.Validation.Email(l.Email)
	c.Validation.Required(l.Username)
	c.Validation.Required(l.Password)
	c.Validation.Required(l.ConfirmPassword)
	c.Validation.MinSize(l.Password, 4).Message("Password must more than 4 characters long")
	c.Validation.Required(l.Password == l.ConfirmPassword).Message("Passwords do not match")
	c.Validation.MinSize(l.ConfirmPassword, 4).Message("Password must more than 4 characters long")
}
