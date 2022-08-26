package services

import (
	"go-booking-api/app/models/responses"
	"strconv"

	"github.com/revel/revel"
)

type ResponseService struct {
}

func NewResponseService() IResponseService {

	return &ResponseService{}
}

func (s *ResponseService) CreateErrorResponse(code int, message string, validationErrors []*revel.ValidationError) *responses.ErrorResponse {

	response := &responses.ErrorResponse{}
	details := []responses.ErrorDetail{}
	for _, error := range validationErrors {
		details = append(details, responses.ErrorDetail{Code: strconv.Itoa(code), Message: error.Message, Target: error.Key})
	}
	response.Error = responses.Error{Code: strconv.Itoa(code), Message: message, Details: details}
	return response
}

func (s *ResponseService) CreateOperationResponse(context string, value interface{}) *responses.OperationResponse {

	response := &responses.OperationResponse{Context: context, Value: value}
	return response
}
