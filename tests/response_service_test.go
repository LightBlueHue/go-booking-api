package tests

import (
	"go-booking-api/app/services"
	"strconv"
	"strings"
	"testing"

	"github.com/revel/revel"
	"github.com/stretchr/testify/assert"
	"syreclabs.com/go/faker"
)

func Test_CreateErrorResponse_ReturnsCorrectly(t *testing.T) {

	target := services.GetResponseService()
	code := faker.RandomInt(100, 599)
	message1 := strings.Join(faker.Lorem().Words(10), " ")
	message2 := strings.Join(faker.Lorem().Words(10), " ")
	message3 := strings.Join(faker.Lorem().Words(10), " ")
	key2 := faker.Lorem().Word()
	key3 := faker.Lorem().Word()
	ve := []*revel.ValidationError{{Message: message2, Key: key2}, {Message: message3, Key: key3}}

	actual := target.CreateErrorResponse(code, message1, ve)

	assert.Equal(t, strconv.Itoa(code), actual.Error.Code)
	assert.Equal(t, message1, actual.Error.Message)

	for i, detail := range actual.Error.Details {
		assert.Equal(t, strconv.Itoa(code), detail.Code)
		assert.Equal(t, ve[i].Message, detail.Message)
		assert.Equal(t, ve[i].Key, detail.Target)
	}
}

func Test_CreateErrorResponse_WithEmptyValidation_ReturnsCorrectly(t *testing.T) {

	target := services.GetResponseService()
	code := faker.RandomInt(100, 599)
	message1 := strings.Join(faker.Lorem().Words(10), " ")
	ve := []*revel.ValidationError{}

	actual := target.CreateErrorResponse(code, message1, ve)

	assert.Equal(t, strconv.Itoa(code), actual.Error.Code)
	assert.Equal(t, message1, actual.Error.Message)
	assert.Empty(t, actual.Error.Details)
}

func Test_CreateOperationResponse_ReturnsCorrectly(t *testing.T) {

	target := services.GetResponseService()
	context := faker.RandomString(5)
	value := faker.RandomString(5)

	actual := target.CreateOperationResponse(context, value)

	assert.Equal(t, context, actual.Context)
	assert.Equal(t, value, actual.Value)
}
