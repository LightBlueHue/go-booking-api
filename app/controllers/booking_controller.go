package controllers

import (
	"fmt"
	"go-booking-api/app/models"
	"go-booking-api/app/models/responses"
	"net/http"

	"github.com/revel/revel"
)

type BookingController struct {
	*revel.Controller
}

func (c *BookingController) Book(count uint) revel.Result {

	_, _, _, rs, _, vs, bs := GetServices()

	vs.ValidateBookingRequest(c.Controller, count)
	if c.Validation.HasErrors() {

		c.Response.Status = http.StatusBadRequest
		resp := rs.CreateErrorResponse(c.Response.Status, "Booking error", c.Validation.Errors)
		return c.RenderJSON(resp)
	}

	var user *models.User
	var err error
	if user, err = GetUser(c.Controller); err != nil {

		c.Validation.Errors = append(c.Validation.Errors, &revel.ValidationError{Message: "unable to authenticate user", Key: "booking"})
		c.Response.Status = http.StatusInternalServerError
		resp := rs.CreateErrorResponse(c.Response.Status, "Booking error", c.Validation.Errors)
		return c.RenderJSON(resp)
	}

	var bookingId uint
	if bookingId, err = bs.Book(user, count); err != nil {

		c.Response.Status = http.StatusInternalServerError
		resp := rs.CreateErrorResponse(c.Response.Status, "Booking error", c.Validation.Errors)
		return c.RenderJSON(resp)
	}

	c.Response.Status = http.StatusCreated
	bookingResponse := responses.BookingResponse{BookingNumber: bookingId, Name: fmt.Sprintf("%s %s", user.FirstName, user.LastName), Tickets: count}
	response := responses.OperationResponse{Value: bookingResponse}
	return c.RenderJSON(response)
}

func (c *BookingController) GetBookings() revel.Result {

	_, _, _, rs, _, _, bs := GetServices()

	var user *models.User
	var bookings *[]models.Booking
	var err error

	if user, err = GetUser(c.Controller); err != nil {

		c.Validation.Errors = append(c.Validation.Errors, &revel.ValidationError{Message: "unable to authenticate user", Key: "booking"})
		c.Response.Status = http.StatusInternalServerError
		resp := rs.CreateErrorResponse(c.Response.Status, "Booking error", c.Validation.Errors)
		return c.RenderJSON(resp)
	}

	if bookings, err = bs.GetBookings(user); err != nil {

		c.Validation.Errors = append(c.Validation.Errors, &revel.ValidationError{Message: "unable to to retrieve bookings", Key: "booking"})
		c.Response.Status = http.StatusInternalServerError
		resp := rs.CreateErrorResponse(c.Response.Status, "Booking retrieval error", c.Validation.Errors)
		return c.RenderJSON(resp)
	}

	c.Response.Status = http.StatusOK
	bookingResponses := []responses.BookingResponse{}

	for _, booking := range *bookings {

		bookingResponses = append(bookingResponses, responses.BookingResponse{BookingNumber: booking.ID, Name: fmt.Sprintf("%s %s", user.FirstName, user.LastName), Tickets: booking.Tickets})
	}

	response := responses.OperationResponse{Value: bookingResponses}
	return c.RenderJSON(response)
}

func init() {

	revel.InterceptFunc(IsLoggedIn, revel.BEFORE, &BookingController{})
}
