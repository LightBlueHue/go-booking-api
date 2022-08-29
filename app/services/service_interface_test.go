package services

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func Test_SetServices_Updates_Correctly(t *testing.T) {

	os.Setenv(GO_BOOKING_API_SECRET, "E59DD115760893782F7FB8CC6C387DE86FFEC3C186A8EFE24184E9CABDB2EFC3")
	var db *gorm.DB

	target := NewService(db, createValidation())

	assert.NotNil(t, target.DBService)
	assert.NotNil(t, target.HashService)
	assert.NotNil(t, target.JWTService)
	assert.NotNil(t, target.ResponseService)
	assert.NotNil(t, target.UserService)
	assert.NotNil(t, target.ValidationService)
	assert.NotNil(t, target.BookingService)
}

func Test_IsServiceSet_WhenAllSet_ReturnsTrue(t *testing.T) {

	os.Setenv(GO_BOOKING_API_SECRET, "E59DD115760893782F7FB8CC6C387DE86FFEC3C186A8EFE24184E9CABDB2EFC3")
	var db *gorm.DB

	target := NewService(db, createValidation())

	var result = target.IsServiceInitialized()

	assert.True(t, result)
}

func Test_IsServiceSet_WhenSomeSet_ReturnsFalse(t *testing.T) {

	os.Setenv(GO_BOOKING_API_SECRET, "E59DD115760893782F7FB8CC6C387DE86FFEC3C186A8EFE24184E9CABDB2EFC3")
	var db *gorm.DB

	target := NewService(db, createValidation())
	target.HashService = nil
	target.ResponseService = nil

	var result = target.IsServiceInitialized()

	assert.False(t, result)
}

func Test_IsServiceSet_WhenAllNotSet_ReturnsFalse(t *testing.T) {

	os.Setenv(GO_BOOKING_API_SECRET, "E59DD115760893782F7FB8CC6C387DE86FFEC3C186A8EFE24184E9CABDB2EFC3")
	var db *gorm.DB

	target := NewService(db, createValidation())
	target.DBService = nil
	target.HashService = nil
	target.JWTService = nil
	target.ResponseService = nil
	target.UserService = nil
	target.ValidationService = nil
	target.BookingService = nil

	var result = target.IsServiceInitialized()

	assert.False(t, result)
}
