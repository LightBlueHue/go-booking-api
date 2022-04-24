package controllers

import (
	"net/http"
	"strconv"

	"github.com/revel/revel"
)

type BookingController struct {
	*revel.Controller
}

func (c BookingController) Book(count uint) revel.Result {

	_, _, _, rs, _, vs, bs := GetServices()

	vs.ValidateBookingRequest(c.Controller, count)
	if c.Validation.HasErrors() {

		c.Response.Status = http.StatusBadRequest
		resp := rs.CreateErrorResponse(c.Response.Status, "Booking error", c.Validation.Errors)
		return c.RenderJSON(resp)
	}

	bs.Book(count)

	return c.RenderText(strconv.FormatUint(uint64(count), 10))
}

func init() {
	revel.InterceptFunc(IsLoggedIn, revel.BEFORE, &BookingController{})
}
