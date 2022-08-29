package controllers

import (
	"fmt"
	"go-booking-api/app/models"
	"go-booking-api/app/models/responses"
	"go-booking-api/app/services"
	"net/http"

	"github.com/revel/revel"
)

// BookingController provides services such as making a Booking and Get bookings by user.
type BookingController struct {
	*revel.Controller
	Service services.Service
}

// Before checks certain conditions are met before endpoints are called. In this case a user has to be logged in.
func (c *BookingController) Before() (result revel.Result, controller *BookingController) {

	return IsLoggedIn(c.Controller, c.Service), c
}

// Book allows the logged in user to book a ticket.
func (c *BookingController) Book(count uint) revel.Result {

	c.Service.ValidationService.ValidateBookingRequest(count)
	if c.Validation.HasErrors() {

		c.Response.Status = http.StatusBadRequest
		resp := c.Service.ResponseService.CreateErrorResponse(c.Response.Status, "Booking error", c.Validation.Errors)
		return c.RenderJSON(resp)
	}

	var user *models.User
	var err error
	if user, err = GetUser(c.Controller, c.Service); err != nil {

		c.Validation.Errors = append(c.Validation.Errors, &revel.ValidationError{Message: "unable to authenticate user", Key: "booking"})
		c.Response.Status = http.StatusInternalServerError
		resp := c.Service.ResponseService.CreateErrorResponse(c.Response.Status, "Booking error", c.Validation.Errors)
		return c.RenderJSON(resp)
	}

	var bookingId uint
	if bookingId, err = c.Service.BookingService.Book(user.ID, count); err != nil {

		c.Response.Status = http.StatusInternalServerError
		resp := c.Service.ResponseService.CreateErrorResponse(c.Response.Status, "Booking error", c.Validation.Errors)
		return c.RenderJSON(resp)
	}

	c.Response.Status = http.StatusCreated
	bookingResponse := responses.BookingResponse{BookingNumber: bookingId, Name: fmt.Sprintf("%s %s", user.FirstName, user.LastName), Tickets: count}
	response := c.Service.ResponseService.CreateOperationResponse("", bookingResponse)
	return c.RenderJSON(response)
}

// GetBookings retrieves booking history of logged in user.
func (c *BookingController) GetBookings() revel.Result {

	var user *models.User
	var bookings *[]models.Booking
	var err error

	if user, err = GetUser(c.Controller, c.Service); err != nil {

		c.Validation.Errors = append(c.Validation.Errors, &revel.ValidationError{Message: "unable to authenticate user", Key: "booking"})
		c.Response.Status = http.StatusInternalServerError
		resp := c.Service.ResponseService.CreateErrorResponse(c.Response.Status, "Booking error", c.Validation.Errors)
		return c.RenderJSON(resp)
	}

	if bookings, err = c.Service.BookingService.GetBookings(user.ID); err != nil {

		c.Validation.Errors = append(c.Validation.Errors, &revel.ValidationError{Message: "unable to to retrieve bookings", Key: "booking"})
		c.Response.Status = http.StatusInternalServerError
		resp := c.Service.ResponseService.CreateErrorResponse(c.Response.Status, "Booking retrieval error", c.Validation.Errors)
		return c.RenderJSON(resp)
	}

	c.Response.Status = http.StatusOK
	bookingResponses := []responses.BookingResponse{}

	for _, booking := range *bookings {

		bookingResponses = append(bookingResponses, responses.BookingResponse{BookingNumber: booking.ID, Name: fmt.Sprintf("%s %s", user.FirstName, user.LastName), Tickets: booking.Tickets})
	}

	response := c.Service.ResponseService.CreateOperationResponse("", bookingResponses)
	return c.RenderJSON(response)
}
