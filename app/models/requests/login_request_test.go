package requests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"syreclabs.com/go/faker"
)

func TestLoginRequest_Behaves_Correctly(t *testing.T) {

	// Arrange
	expectedEmail:= faker.Internet().Email()
	expectedPassword:=faker.Internet().Password(4,20)
	
	// Act
	target:= &LoginRequest{Email: expectedEmail,Password: expectedPassword}

	// Assert
	assert.Equal(t,expectedEmail, target.Email)
	assert.Equal(t,expectedPassword, target.Password)
}