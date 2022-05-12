package tests

import (
	"go-booking-api/app/models/requests"
	"go-booking-api/app/services"
	"testing"

	"github.com/revel/revel"
	"github.com/stretchr/testify/assert"
	"syreclabs.com/go/faker"
)

func Test_ValidateLoginRequest_InValidEmail_Returns_Error(t *testing.T) {

	target := services.GetValidationService()
	rv := createValidation()
	requests := []requests.LoginRequest{
		{Email: "faker.Internet().Email()", Password: faker.Internet().Password(4, 20)},
		{Email: "", Password: faker.Internet().Password(4, 20)},
		{Email: "ya.com", Password: faker.Internet().Password(4, 20)},
	}

	for i, request := range requests {

		target.ValidateLoginRequest(rv, &request)

		assert.True(t, rv.HasErrors())
		assert.Contains(t, rv.Errors[i].Message, "email")
	}
}

func Test_ValidateLoginRequest_ValidEmail_Returns_NoError(t *testing.T) {

	target := services.GetValidationService()
	rv := createValidation()
	requests := []requests.LoginRequest{
		{Email: faker.Internet().Email(), Password: faker.Internet().Password(4, 20)},
		{Email: faker.Internet().FreeEmail(), Password: faker.Internet().Password(4, 20)},
		{Email: faker.Internet().SafeEmail(), Password: faker.Internet().Password(4, 20)},
	}

	for _, request := range requests {

		target.ValidateLoginRequest(rv, &request)

		assert.False(t, rv.HasErrors())
	}
}

func Test_ValidateLoginRequest_InValidPassword_Returns_Error(t *testing.T) {

	target := services.GetValidationService()
	rv := createValidation()
	requests := []requests.LoginRequest{
		{Email: faker.Internet().Email(), Password: faker.Internet().Password(0, 0)},
		{Email: faker.Internet().Email(), Password: faker.Internet().Password(0, 3)},
		{Email: faker.Internet().FreeEmail(), Password: faker.Internet().Password(20, 30)},
	}

	for i, request := range requests {

		target.ValidateLoginRequest(rv, &request)

		assert.True(t, rv.HasErrors())
		assert.Equal(t, services.VALIDATION_REQUEST_PASSWORD_LENGTH, rv.Errors[i].Message)
	}
}

func Test_ValidateLoginRequest_ValidPassword_Returns_NoError(t *testing.T) {

	target := services.GetValidationService()
	rv := createValidation()
	requests := []requests.LoginRequest{
		{Email: faker.Internet().Email(), Password: faker.Internet().Password(4, 4)},
		{Email: faker.Internet().Email(), Password: faker.Internet().Password(4, 20)},
		{Email: faker.Internet().FreeEmail(), Password: faker.Internet().Password(5, 15)},
		{Email: faker.Internet().FreeEmail(), Password: faker.Internet().Password(20, 20)},
	}

	for _, request := range requests {

		target.ValidateLoginRequest(rv, &request)

		assert.False(t, rv.HasErrors())
	}
}

func Test_ValidateRegisterRequest_InValidFirstName_Returns_Error(t *testing.T) {

	target := services.GetValidationService()
	rv := createValidation()
	pwd := faker.Internet().Password(4, 20)
	request := requests.RegisterRequest{
		FirstName:       "",
		LastName:        faker.Name().LastName(),
		Email:           faker.Internet().Email(),
		Password:        pwd,
		ConfirmPassword: pwd,
	}

	target.ValidateRegisterRequest(rv, &request)

	assert.True(t, rv.HasErrors())
	assert.Equal(t, services.VALIDATION_REGISTER_REQUEST_FIRSTNAME_INVALID, rv.Errors[0].Message)
}

func Test_ValidateRegisterRequest_InValidLastName_Returns_Error(t *testing.T) {

	target := services.GetValidationService()
	rv := createValidation()
	pwd := faker.Internet().Password(4, 20)
	request := requests.RegisterRequest{
		FirstName:       faker.Name().FirstName(),
		LastName:        "",
		Email:           faker.Internet().Email(),
		Password:        pwd,
		ConfirmPassword: pwd,
	}

	target.ValidateRegisterRequest(rv, &request)

	assert.True(t, rv.HasErrors())
	assert.Equal(t, services.VALIDATION_REGISTER_REQUEST_LASTNAME_INVALID, rv.Errors[0].Message)
}

func Test_ValidateRegisterRequest_InValidEmail_Returns_Error(t *testing.T) {

	target := services.GetValidationService()
	rv := createValidation()
	pwd := faker.Internet().Password(4, 20)
	firstName := faker.Name().FirstName()
	lastName := faker.Name().LastName()
	requests := []requests.RegisterRequest{
		{FirstName: firstName, LastName: lastName, Email: "faker.Internet().Email()", Password: pwd, ConfirmPassword: pwd},
		{FirstName: firstName, LastName: lastName, Email: "", Password: pwd, ConfirmPassword: pwd},
		{FirstName: firstName, LastName: lastName, Email: "ya.com", Password: pwd, ConfirmPassword: pwd},
	}

	for i, request := range requests {

		target.ValidateRegisterRequest(rv, &request)

		assert.True(t, rv.HasErrors())
		assert.Contains(t, rv.Errors[i].Message, "email")
	}
}

func Test_ValidateRegisterRequest_InValidPassword_Returns_Error(t *testing.T) {

	target := services.GetValidationService()
	rv := createValidation()
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

	for i, request := range requests {

		target.ValidateRegisterRequest(rv, &request)

		assert.True(t, rv.HasErrors())
		assert.Equal(t, services.VALIDATION_REQUEST_PASSWORD_LENGTH, rv.Errors[i].Message)
	}
}

func Test_ValidateRegisterRequest_NonMatchingPasswords_Returns_Error(t *testing.T) {

	target := services.GetValidationService()
	rv := createValidation()
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

	for i, request := range requests {

		target.ValidateRegisterRequest(rv, &request)

		assert.True(t, rv.HasErrors())
		assert.Equal(t, services.VALIDATION_REGISTER_REQUEST_PASSWORD_NOMATCH, rv.Errors[i].Message)
	}
}

func Test_ValidateBookingRequest_InValidNumber_Returns_Error(t *testing.T) {

	target := services.GetValidationService()
	rv := createValidation()

	target.ValidateBookingRequest(rv, 0)

	assert.True(t, rv.HasErrors())
	assert.Equal(t, services.VALIDATION_BOOKING_REQUEST_VALID_NUMBER, rv.Errors[0].Message)
}

func Test_ValidateBookingRequest_ValidNumber_Returns_NoError(t *testing.T) {

	target := services.GetValidationService()
	rv := createValidation()

	target.ValidateBookingRequest(rv, uint(faker.RandomInt(1, 1000)))

	assert.False(t, rv.HasErrors())
}

func createValidation() *revel.Validation {

	c := revel.NewControllerEmpty()
	c.Validation = &revel.Validation{Request: c.Request, Translator: revel.MessageFunc}
	return c.Validation
}
