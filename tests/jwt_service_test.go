package tests

import (
	"go-booking-api/app/services"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"syreclabs.com/go/faker"
)

func Test_GenerateToken_Returns_Token(t *testing.T) {

	data := make(map[string]bool)
	data[faker.Internet().Email()] = false
	data[faker.Internet().Email()] = true
	data[""] = true
	target := services.GetJWTService()

	for email, isUser := range data {

		token := target.GenerateToken(email, isUser)

		assert.NotEmpty(t, token)
	}
}

func Test_ValidateToken_ValidToken_Returns_Token(t *testing.T) {

	target := services.GetJWTService()
	email := faker.Internet().Email()
	isUser := true
	token := target.GenerateToken(email, isUser)

	actual, err := target.ValidateToken(token)

	assert.Nil(t, err)
	assert.NotNil(t, actual)
}

func Test_ValidateToken_WrongTokenFormat_Returns_Error(t *testing.T) {

	os.Setenv(services.GO_BOOKING_API_SECRET, "E59DD115760893782F7FB8CC6C387DE86FFEC3C186A8EFE24184E9CABDB2EFC3")
	token := faker.RandomString(229)
	target := services.GetJWTService()

	actual, err := target.ValidateToken(token)

	assert.Error(t, err)
	assert.Nil(t, actual)
}

func Test_ValidateToken_ExpiredToken_Returns_Error(t *testing.T) {

	os.Setenv(services.GO_BOOKING_API_SECRET, "E59DD115760893782F7FB8CC6C387DE86FFEC3C186A8EFE24184E9CABDB2EFC3")
	expiredToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoieWFjQHlhYy5jbyIsImVtYWlsIjoieWFjQHlhYy5jbyIsInVzZXIiOnRydWUsImV4cCI6MTY1MjAzMzc1MywiaWF0IjoxNjUyMDMzNzUyLCJpc3MiOiJnby1ib29raW5nLWFwaSJ9.-t_MwxN_sQ5pxG6B0X1zrtdCnefCcXfvi1byCrpDQjg"
	target := services.GetJWTService()

	actual, err := target.ValidateToken(expiredToken)

	assert.Error(t, err)
	assert.NotNil(t, actual)
	assert.False(t, actual.Valid)
}

func Test_GetClaim_ValidToken_Returns_CorrectClaim(t *testing.T) {

	os.Setenv(services.GO_BOOKING_API_SECRET, "E59DD115760893782F7FB8CC6C387DE86FFEC3C186A8EFE24184E9CABDB2EFC3")
	target := services.GetJWTService()
	email := faker.Internet().Email()
	isUser := true
	token := target.GenerateToken(email, isUser)

	actual, err := target.GetClaim(token, services.EMAIL_CLAIM)

	assert.Nil(t, err)
	assert.Equal(t, email, actual)
}

func Test_GetClaim_ExpiredToken_Returns_Error(t *testing.T) {

	os.Setenv(services.GO_BOOKING_API_SECRET, "E59DD115760893782F7FB8CC6C387DE86FFEC3C186A8EFE24184E9CABDB2EFC3")
	expiredToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoieWFjQHlhYy5jbyIsImVtYWlsIjoieWFjQHlhYy5jbyIsInVzZXIiOnRydWUsImV4cCI6MTY1MjAzMzc1MywiaWF0IjoxNjUyMDMzNzUyLCJpc3MiOiJnby1ib29raW5nLWFwaSJ9.-t_MwxN_sQ5pxG6B0X1zrtdCnefCcXfvi1byCrpDQjg"
	target := services.GetJWTService()

	actual, err := target.GetClaim(expiredToken, services.EMAIL_CLAIM)

	assert.Error(t, err)
	assert.Empty(t, actual)
}

func Test_GetClaim_ValidToken_WithInValidClaimRequest_Returns_Error(t *testing.T) {

	os.Setenv(services.GO_BOOKING_API_SECRET, "E59DD115760893782F7FB8CC6C387DE86FFEC3C186A8EFE24184E9CABDB2EFC3")
	target := services.GetJWTService()
	email := faker.Internet().Email()
	isUser := true
	token := target.GenerateToken(email, isUser)
	claims := []string{"services.EMAIL_CLAIM", ""}

	for _, claim := range claims {

		actual, err := target.GetClaim(token, services.JwtClaimType(claim))

		assert.Error(t, err)
		assert.Empty(t, actual)
	}
}
