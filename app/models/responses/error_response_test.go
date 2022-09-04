package responses_test

import (
	"go-booking-api/app/models/responses"
	"testing"

	"github.com/stretchr/testify/assert"
	"syreclabs.com/go/faker"
)

func TestErrorResponse_Behaves_Correctly(t *testing.T) {

	// Arrange
	expectedErrorCode := faker.RandomString(5)
	expectedErrorMessage := faker.Lorem().Sentence(10)
	expectedErrorTarget := faker.RandomString(5)

	expectedErrorDetailCode := faker.RandomString(5)
	expectedErrorDetailMessage := faker.Lorem().Sentence(10)
	expectedErrorDetailTarget := faker.RandomString(5)

	// Act
	actual := responses.ErrorResponse{
		Error: responses.Error{
			Code:    expectedErrorCode,
			Message: expectedErrorMessage,
			Target:  expectedErrorTarget,
			Details: []responses.ErrorDetail{
				{
					Code:    expectedErrorDetailCode,
					Message: expectedErrorDetailMessage,
					Target:  expectedErrorDetailTarget,
				},
			},
		}}

	// Assert
	assert.Equal(t, expectedErrorCode, actual.Error.Code)
	assert.Equal(t, expectedErrorMessage, actual.Error.Message)
	assert.Equal(t, expectedErrorTarget, actual.Error.Target)

	assert.Equal(t, expectedErrorDetailCode, actual.Error.Details[0].Code)
	assert.Equal(t, expectedErrorDetailMessage, actual.Error.Details[0].Message)
	assert.Equal(t, expectedErrorDetailTarget, actual.Error.Details[0].Target)
}
