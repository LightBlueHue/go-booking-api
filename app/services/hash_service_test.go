package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"syreclabs.com/go/faker"
)

func Test_HashAndSalt_ValidPassword_Returns_HashedPassword(t *testing.T) {

	// Arrange
	target := NewHashService()
	pwd := faker.Internet().Password(4, 20)

	// Act
	actual, err := target.HashAndSalt(pwd)

	// Assert
	assert.Nil(t, err)
	assert.NotNil(t, actual)
	assert.NotEmpty(t, actual)
}

func Test_CompareHashAndPassword_WithCorrectData_Returns_True(t *testing.T) {

	// Arrange
	target := NewHashService()

	pwd := faker.Internet().Password(4, 20)
	hshPwd, err := target.HashAndSalt(pwd)

	assert.Nil(t, err)

	// Act
	actual, err := target.CompareHashAndPassword(hshPwd, pwd)

	// Assert
	assert.True(t, actual)
	assert.Nil(t, err)
}

func Test_CompareHashAndPassword_WithInValidData_Returns_False(t *testing.T) {

	// Arrange
	pwds := []string{faker.Internet().Password(4, 20), ""}
	target := NewHashService()

	// Act
	hshPwd, err := target.HashAndSalt(faker.Internet().Password(4, 20))

	// Assert
	for _, pwd := range pwds {

		assert.Nil(t, err)

		actual, err := target.CompareHashAndPassword(hshPwd, pwd)

		assert.False(t, actual)
		assert.NotNil(t, err)
	}
}
