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

	_, _, _, rs, us, vs, bs := GetServices()

	vs.ValidateBookingRequest(c.Controller, count)
	if c.Validation.HasErrors() {

		c.Response.Status = http.StatusBadRequest
		resp := rs.CreateErrorResponse(c.Response.Status, "Booking error", c.Validation.Errors)
		return c.RenderJSON(resp)
	}

	var user *models.User
	var err error
	token, _ := GetBearerToken(c.Controller)
	if user, err = us.GetByToken(token); err != nil {

		c.Validation.Errors = append(c.Validation.Errors, &revel.ValidationError{Message: "token", Key: "booking"})
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

func init() {
	
	revel.InterceptFunc(IsLoggedIn, revel.BEFORE, &BookingController{})
}
