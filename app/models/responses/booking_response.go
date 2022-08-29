// Package responses contains api response models
package responses

// BookingResponse contains information about the bookings a user made.
type BookingResponse struct {
	BookingNumber uint
	Name          string
	Tickets       uint
}
