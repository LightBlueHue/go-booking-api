package tests

import (
	"go-booking-api/app/services"
	"testing"

	"github.com/stretchr/testify/assert"
	"syreclabs.com/go/faker"
)

func Test_HashAndSalt_ValidPassword_ReturnsCorrectData(t *testing.T) {

	target := services.GetHashService()

	pwd := faker.Internet().Password(4, 20)
	actual, err := target.HashAndSalt(pwd)

	assert.Nil(t, err)
	assert.NotNil(t, actual)
	assert.Greater(t, len(actual), 0)
}

func Test_ComparePasswords_WithCorrectData_ReturnsTrue(t *testing.T) {

	target := services.GetHashService()

	pwd := faker.Internet().Password(4, 20)
	hshPwd, err := target.HashAndSalt(pwd)

	assert.Nil(t, err)

	actual := target.ComparePasswords(hshPwd, pwd)

	assert.True(t, actual)
}

func Test_ComparePasswords_WithInValidData_ReturnsFalse(t *testing.T) {

	pwds := []string{faker.Internet().Password(4, 20), ""}
	target := services.GetHashService()

	hshPwd, err := target.HashAndSalt(faker.Internet().Password(4, 20))

	for _, pwd := range pwds {

		assert.Nil(t, err)

		actual := target.ComparePasswords(hshPwd, pwd)

		assert.False(t, actual)
	}
}
