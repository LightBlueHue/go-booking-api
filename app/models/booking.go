// Package models contains database models.
package models

import "gorm.io/gorm"

// Booking represents a booking by a user.
type Booking struct {
	gorm.Model
	UserID            uint
	User              User
	TicketInventoryID uint
	TicketInventory   TicketInventory
	Tickets           uint
}
