package models

type Booking struct {
	*User
	Tickets uint
}