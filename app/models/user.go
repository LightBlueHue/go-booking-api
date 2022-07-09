package models

import "gorm.io/gorm"

// User represents and holds information about a user.
type User struct {
	gorm.Model
	FirstName    string
	LastName     string
	Email        string `gorm:"uniqueIndex;size:20"`
	CredentialID uint
	Credential   Credentials
	Booking      []Booking
}
