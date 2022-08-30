package services_test

import (
	"go-booking-api/app/services"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"syreclabs.com/go/faker"
)

func Test_GenerateToken_Returns_Token(t *testing.T) {

	// Arrange
	data := make(map[string]bool)
	data[faker.Internet().Email()] = false
	data[faker.Internet().Email()] = true
	data[""] = true
	target := services.NewJWTService()

	for email, isUser := range data {

		// Act
		token := target.GenerateToken(email, isUser)

		// Assert
		assert.NotEmpty(t, token)
	}
}

func Test_ValidateToken_WhenTokenValid_Returns_Token(t *testing.T) {

	// Arrange
	target := services.NewJWTService()
	email := faker.Internet().Email()
	isUser := true
	token := target.GenerateToken(email, isUser)

	// Act
	actual, err := target.ValidateToken(token)

	// Assert
	assert.Nil(t, err)
	assert.NotNil(t, actual)
}

func Test_ValidateToken_WhenWrongTokenFormat_Returns_Error(t *testing.T) {

	// Arrange
	os.Setenv(services.GO_BOOKING_API_SECRET, "E59DD115760893782F7FB8CC6C387DE86FFEC3C186A8EFE24184E9CABDB2EFC3")
	token := faker.RandomString(229)
	target := services.NewJWTService()

	// Act
	actual, err := target.ValidateToken(token)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, actual)
}

func Test_ValidateToken_WhenTokenExpired_Returns_Error(t *testing.T) {

	// Arrange
	os.Setenv(services.GO_BOOKING_API_SECRET, "E59DD115760893782F7FB8CC6C387DE86FFEC3C186A8EFE24184E9CABDB2EFC3")
	expiredToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoieWFjQHlhYy5jbyIsImVtYWlsIjoieWFjQHlhYy5jbyIsInVzZXIiOnRydWUsImV4cCI6MTY1MjAzMzc1MywiaWF0IjoxNjUyMDMzNzUyLCJpc3MiOiJnby1ib29raW5nLWFwaSJ9.-t_MwxN_sQ5pxG6B0X1zrtdCnefCcXfvi1byCrpDQjg"
	target := services.NewJWTService()

	// Act
	actual, err := target.ValidateToken(expiredToken)

	// Assert
	assert.Error(t, err)
	assert.NotNil(t, actual)
	assert.False(t, actual.Valid)
}

func Test_GetClaim_WhenTokenValid_Returns_CorrectClaim(t *testing.T) {

	// Arrange
	os.Setenv(services.GO_BOOKING_API_SECRET, "E59DD115760893782F7FB8CC6C387DE86FFEC3C186A8EFE24184E9CABDB2EFC3")
	target := services.NewJWTService()
	email := faker.Internet().Email()
	isUser := true
	token := target.GenerateToken(email, isUser)

	// Act
	actual, err := target.GetClaim(token, services.EMAIL_CLAIM)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, email, actual)
}

func Test_GetClaim_WhenTokenExpired_Returns_Error(t *testing.T) {

	// Arrange
	os.Setenv(services.GO_BOOKING_API_SECRET, "E59DD115760893782F7FB8CC6C387DE86FFEC3C186A8EFE24184E9CABDB2EFC3")
	expiredToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoieWFjQHlhYy5jbyIsImVtYWlsIjoieWFjQHlhYy5jbyIsInVzZXIiOnRydWUsImV4cCI6MTY1MjAzMzc1MywiaWF0IjoxNjUyMDMzNzUyLCJpc3MiOiJnby1ib29raW5nLWFwaSJ9.-t_MwxN_sQ5pxG6B0X1zrtdCnefCcXfvi1byCrpDQjg"
	target := services.NewJWTService()

	// Act
	actual, err := target.GetClaim(expiredToken, services.EMAIL_CLAIM)

	// Assert
	assert.Error(t, err)
	assert.Empty(t, actual)
}

func Test_GetClaim_WhenTokenValid_ButWrongClaim_Returns_Error(t *testing.T) {

	// Arrange
	os.Setenv(services.GO_BOOKING_API_SECRET, "E59DD115760893782F7FB8CC6C387DE86FFEC3C186A8EFE24184E9CABDB2EFC3")
	target := services.NewJWTService()
	email := faker.Internet().Email()
	isUser := true
	token := target.GenerateToken(email, isUser)
	claims := []string{faker.RandomString(7), ""}

	for _, claim := range claims {

		// Act
		actual, err := target.GetClaim(token, services.JwtClaimType(claim))

		// Asserts
		assert.Error(t, err)
		assert.Empty(t, actual)
	}
}

func Test_GetSecretKey_WhenSecretSetInEnvironmentVariable_Returns_Secret(t *testing.T) {

	// Arrange
	expectedSecrets := []string{"E59DD115760893782F7FB8CC6C387DE86FFEC3C186A8EFE24184E9CABDB2EFC3", " "}

	for _, expectedSecret := range expectedSecrets {

		os.Setenv(services.GO_BOOKING_API_SECRET, expectedSecret)

		// Act
		actualSecret := services.GetSecretKey()

		// Assert
		assert.Equal(t, expectedSecret, actualSecret)
	}
}

func Test_GetSecretKey_WhenEnvironmentVariableEmpty_Panics(t *testing.T) {

	// Arrange
	os.Setenv(services.GO_BOOKING_API_SECRET, "")

	// Act
	// Assert
	assert.PanicsWithValue(t, services.GO_BOOKING_API_SECRET, func() { services.GetSecretKey() })
}
