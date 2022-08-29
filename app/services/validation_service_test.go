package services

import (
	"go-booking-api/app/models/requests"
	"os"
	"testing"

	"github.com/revel/revel"
	"github.com/stretchr/testify/assert"
	"syreclabs.com/go/faker"
)

func Test_ValidateLoginRequest_WhenEmailInValid_Returns_Error(t *testing.T) {

	// Arrange
	target := NewValidationService()
	requests := []requests.LoginRequest{
		{Email: "faker.Internet().Email()", Password: faker.Internet().Password(4, 20)},
		{Email: "", Password: faker.Internet().Password(4, 20)},
		{Email: "ya.com", Password: faker.Internet().Password(4, 20)},
	}

	for _, request := range requests {

		rv := createValidation()

		// Act
		target.ValidateLoginRequest(rv, &request)

		// Assert
		assert.True(t, rv.HasErrors())
		assert.Contains(t, rv.Errors[0].Message, "email")
	}
}

func Test_ValidateLoginRequest_WhenEmailValid_Returns_NoError(t *testing.T) {

	// Arrange
	target := NewValidationService()
	requests := []requests.LoginRequest{
		{Email: faker.Internet().Email(), Password: faker.Internet().Password(4, 20)},
		{Email: faker.Internet().FreeEmail(), Password: faker.Internet().Password(4, 20)},
		{Email: faker.Internet().SafeEmail(), Password: faker.Internet().Password(4, 20)},
	}

	for _, request := range requests {

		rv := createValidation()

		// Act
		target.ValidateLoginRequest(rv, &request)

		// Assert
		assert.False(t, rv.HasErrors())
	}
}

func Test_ValidateLoginRequest_WhenPasswordInValid_Returns_Error(t *testing.T) {

	// Arrange
	target := NewValidationService()
	requests := []requests.LoginRequest{
		{Email: faker.Internet().Email(), Password: faker.Internet().Password(0, 0)},
		{Email: faker.Internet().Email(), Password: faker.Internet().Password(0, 3)},
		{Email: faker.Internet().FreeEmail(), Password: faker.Internet().Password(20, 30)},
	}

	for _, request := range requests {

		rv := createValidation()
		// Act
		target.ValidateLoginRequest(rv, &request)

		// Assert
		assert.True(t, rv.HasErrors())
		assert.Equal(t, VALIDATION_REQUEST_PASSWORD_LENGTH, rv.Errors[0].Message)
	}
}

func Test_ValidateLoginRequest_WhenPasswordValid_Returns_NoError(t *testing.T) {

	// Arrange
	target := NewValidationService()
	requests := []requests.LoginRequest{
		{Email: faker.Internet().Email(), Password: faker.Internet().Password(4, 4)},
		{Email: faker.Internet().Email(), Password: faker.Internet().Password(4, 20)},
		{Email: faker.Internet().FreeEmail(), Password: faker.Internet().Password(5, 15)},
		{Email: faker.Internet().FreeEmail(), Password: faker.Internet().Password(20, 20)},
	}

	for _, request := range requests {

		rv := createValidation()

		// Act
		target.ValidateLoginRequest(rv, &request)

		// Assert
		assert.False(t, rv.HasErrors())
	}
}

func Test_ValidateRegisterRequest_WhenFirstNameInValid_Returns_Error(t *testing.T) {

	// Arrange
	target := NewValidationService()
	rv := createValidation()
	pwd := faker.Internet().Password(4, 20)
	request := requests.RegisterRequest{
		FirstName:       "",
		LastName:        faker.Name().LastName(),
		Email:           faker.Internet().Email(),
		Password:        pwd,
		ConfirmPassword: pwd,
	}

	// Act
	target.ValidateRegisterRequest(rv, &request)

	// Assert
	assert.True(t, rv.HasErrors())
	assert.Equal(t, VALIDATION_REGISTER_REQUEST_FIRSTNAME_INVALID, rv.Errors[0].Message)
}

func Test_ValidateRegisterRequest_WhenLastNameInValid_Returns_Error(t *testing.T) {

	// Arrange
	target := NewValidationService()
	rv := createValidation()
	pwd := faker.Internet().Password(4, 20)
	request := requests.RegisterRequest{
		FirstName:       faker.Name().FirstName(),
		LastName:        "",
		Email:           faker.Internet().Email(),
		Password:        pwd,
		ConfirmPassword: pwd,
	}

	// Act
	target.ValidateRegisterRequest(rv, &request)

	// Assert
	assert.True(t, rv.HasErrors())
	assert.Equal(t, VALIDATION_REGISTER_REQUEST_LASTNAME_INVALID, rv.Errors[0].Message)
}

func Test_ValidateRegisterRequest_WhenEmailInValid_Returns_Error(t *testing.T) {

	// Arrange
	target := NewValidationService()
	pwd := faker.Internet().Password(4, 20)
	firstName := faker.Name().FirstName()
	lastName := faker.Name().LastName()
	requests := []requests.RegisterRequest{
		{FirstName: firstName, LastName: lastName, Email: "faker.Internet().Email()", Password: pwd, ConfirmPassword: pwd},
		{FirstName: firstName, LastName: lastName, Email: "", Password: pwd, ConfirmPassword: pwd},
		{FirstName: firstName, LastName: lastName, Email: "ya.com", Password: pwd, ConfirmPassword: pwd},
	}

	for _, request := range requests {

		rv := createValidation()

		// Act
		target.ValidateRegisterRequest(rv, &request)

		// Assert
		assert.True(t, rv.HasErrors())
		assert.Contains(t, rv.Errors[0].Message, "email")
	}
}

func Test_ValidateRegisterRequest_WhenPasswordInValid_Returns_Error(t *testing.T) {

	// Arrange
	os.Setenv(GO_BOOKING_API_SECRET, "E59DD115760893782F7FB8CC6C387DE86FFEC3C186A8EFE24184E9CABDB2EFC3")
	target := NewValidationService()
	email := faker.Internet().Email()
	pwd1 := faker.Internet().Password(0, 0)
	pwd2 := faker.Internet().Password(0, 3)
	pwd3 := faker.Internet().Password(20, 30)
	firstName := faker.Name().FirstName()
	lastName := faker.Name().LastName()

	requests := []requests.RegisterRequest{
		{FirstName: firstName, LastName: lastName, Email: email, Password: pwd1, ConfirmPassword: pwd1},
		{FirstName: firstName, LastName: lastName, Email: email, Password: pwd2, ConfirmPassword: pwd2},
		{FirstName: firstName, LastName: lastName, Email: email, Password: pwd3, ConfirmPassword: pwd3},
	}

	for _, request := range requests {

		rv := createValidation()

		// Act
		target.ValidateRegisterRequest(rv, &request)

		// Assert
		assert.True(t, rv.HasErrors())
		assert.Equal(t, VALIDATION_REQUEST_PASSWORD_LENGTH, rv.Errors[0].Message)
	}
}

func Test_ValidateRegisterRequest_WhenPasswordsDoNotMatch_Returns_Error(t *testing.T) {

	// Arrange
	target := NewValidationService()
	email := faker.Internet().Email()
	pwd := faker.Internet().Password(4, 6)
	confPwd1 := faker.Internet().Password(0, 0)
	confPwd2 := faker.Internet().Password(0, 3)
	confPwd3 := faker.Internet().Password(20, 30)
	firstName := faker.Name().FirstName()
	lastName := faker.Name().LastName()

	requests := []requests.RegisterRequest{
		{FirstName: firstName, LastName: lastName, Email: email, Password: pwd, ConfirmPassword: confPwd1},
		{FirstName: firstName, LastName: lastName, Email: email, Password: pwd, ConfirmPassword: confPwd2},
		{FirstName: firstName, LastName: lastName, Email: email, Password: pwd, ConfirmPassword: confPwd3},
	}

	for _, request := range requests {

		rv := createValidation()

		// Act
		target.ValidateRegisterRequest(rv, &request)

		// Assert
		assert.True(t, rv.HasErrors())
		assert.Equal(t, VALIDATION_REGISTER_REQUEST_PASSWORD_NOMATCH, rv.Errors[0].Message)
	}
}

func Test_ValidateBookingRequest_WhenNumberInValid_Returns_Error(t *testing.T) {

	// Arrange
	target := NewValidationService()
	rv := createValidation()

	// Act
	target.ValidateBookingRequest(rv, 0)

	// Assert
	assert.True(t, rv.HasErrors())
	assert.Equal(t, VALIDATION_BOOKING_REQUEST_VALID_NUMBER, rv.Errors[0].Message)
}

func Test_ValidateBookingRequest_WhenNumberValid_Returns_NoError(t *testing.T) {

	// Arrange
	target := NewValidationService()
	rv := createValidation()

	// Act
	target.ValidateBookingRequest(rv, uint(faker.RandomInt(1, 1000)))

	// Assert
	assert.False(t, rv.HasErrors())
}

func createValidation() *revel.Validation {

	c := revel.NewControllerEmpty()
	c.Validation = &revel.Validation{Request: c.Request, Translator: revel.MessageFunc}
	return c.Validation
}
