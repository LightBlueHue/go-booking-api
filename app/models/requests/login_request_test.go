package requests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"syreclabs.com/go/faker"
)

func TestLoginRequest_Behaves_Correctly(t *testing.T) {

	// Arrange
	expectedEmail := faker.Internet().Email()
	expectedPassword := faker.Internet().Password(4, 20)

	// Act
	actual := &LoginRequest{Email: expectedEmail, Password: expectedPassword}

	// Assert
	assert.Equal(t, expectedEmail, actual.Email)
	assert.Equal(t, expectedPassword, actual.Password)
}
