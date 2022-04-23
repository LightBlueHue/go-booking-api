package models

import "gorm.io/gorm"

type Booking struct {
	gorm.Model
	UserID  uint
	User    User
	Tickets uint
}
