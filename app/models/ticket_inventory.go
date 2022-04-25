package models

import (
	"gorm.io/gorm"
)

type TicketInventory struct {
	gorm.Model
	AvailableTickets uint
	TotalTickets     uint
	Name             string `gorm:"size:256"`
	Description      string `gorm:"size:256"`
}
