package requests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"syreclabs.com/go/faker"
)

func TestRegisterRequest_Behaves_Correctly(t *testing.T) {

	// Arrange
	expectedFirstName := faker.Name().FirstName()
	expectedLastName := faker.Name().LastName()
	expectedEmail := faker.Internet().Email()
	expectedPassword := faker.Internet().Password(4, 20)
	expectedConfirmPassword := faker.Internet().Password(4, 20)

	// Act
	actual := RegisterRequest{FirstName: expectedFirstName, LastName: expectedLastName, Email: expectedEmail, Password: expectedPassword, ConfirmPassword: expectedConfirmPassword}

	// Assert
	assert.Equal(t, expectedFirstName, actual.FirstName)
	assert.Equal(t, expectedLastName, actual.LastName)
	assert.Equal(t, expectedEmail, actual.Email)
	assert.Equal(t, expectedPassword, actual.Password)
	assert.Equal(t, expectedConfirmPassword, actual.ConfirmPassword)
}
