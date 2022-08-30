package responses_test

import (
	"go-booking-api/app/models/responses"
	"testing"

	"github.com/stretchr/testify/assert"
	"syreclabs.com/go/faker"
)

func TestOperationResponse__Behaves_Correctly(t *testing.T) {

	// Arrange
	var expectedContext string
	expectedValues := make([]interface{}, 3)
	expectedValues = append(expectedValues, faker.Lorem().Sentence(10))
	expectedValues = append(expectedValues, "")
	expectedValues = append(expectedValues, responses.BookingResponse{BookingNumber: uint(faker.RandomInt(1, 100)), Name: faker.Name().Name(), Tickets: uint(faker.RandomInt(1, 100))})

	for i, expectedValue := range expectedValues {

		expectedContext = faker.Lorem().Sentence(10)

		// Act
		actual := responses.OperationResponse{Context: expectedContext, Value: expectedValue}

		// Assert
		assert.Equal(t, expectedContext, actual.Context)
		assert.Equal(t, expectedValue, expectedValues[i])
	}
}
