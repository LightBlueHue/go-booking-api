package responses_test

import (
	"go-booking-api/app/models/responses"
	"testing"

	"github.com/stretchr/testify/assert"
	"syreclabs.com/go/faker"
)

func TestBookingResponse_Behaves_Correctly(t *testing.T) {

	// Arrange
	expectedBookingNumber := uint(faker.RandomInt(0, 100))
	expectedName := faker.Name().Name()
	expectedTickets := uint(faker.RandomInt(0, 100))

	// Act
	actual := responses.BookingResponse{BookingNumber: uint(expectedBookingNumber), Name: expectedName, Tickets: uint(expectedTickets)}

	// Assert
	assert.Equal(t, expectedBookingNumber, actual.BookingNumber)
	assert.Equal(t, expectedName, actual.Name)
	assert.Equal(t, expectedTickets, actual.Tickets)
}
