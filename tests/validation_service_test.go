package tests

import (
	"go-booking-api/app/models/requests"
	"go-booking-api/app/services"
	"testing"

	"github.com/revel/revel"
	"github.com/stretchr/testify/assert"
	"syreclabs.com/go/faker"
)

func Test_ValidateLoginRequest_InValidEmail_ReturnsError(t *testing.T) {

	target := services.GetValidationService()
	c := revel.NewControllerEmpty()
	c.Validation = &revel.Validation{Request: c.Request, Translator: revel.MessageFunc}
	rv := c.Validation
	request := requests.LoginRequest{Email: "xx.com", Password: faker.Internet().Password(4, 20)}

	target.ValidateLoginRequest(rv, &request)

	assert.Equal(t, true, rv.HasErrors())
	assert.Contains(t, rv.Errors[0].Message, "email")
}

func Test_ValidateLoginRequest_ValidEmail_ReturnsNoError(t *testing.T) {

	requests := []requests.LoginRequest{
		{Email: faker.Internet().Email(), Password: faker.Internet().Password(4, 20)},
		{Email: faker.Internet().FreeEmail(), Password: faker.Internet().Password(4, 20)},
		{Email: faker.Internet().SafeEmail(), Password: faker.Internet().Password(4, 20)},
	}
	target := services.GetValidationService()
	c := revel.NewControllerEmpty()
	c.Validation = &revel.Validation{Request: c.Request, Translator: revel.MessageFunc}
	rv := c.Validation

	for _, request := range requests {

		target.ValidateLoginRequest(rv, &request)
		assert.Equal(t, false, rv.HasErrors())

	}

}
