package models

import "gorm.io/gorm"

type Booking struct {
	gorm.Model
	UserID            uint
	User              User
	TicketInventoryID uint
	TicketInventory   TicketInventory
	Tickets           uint
}
