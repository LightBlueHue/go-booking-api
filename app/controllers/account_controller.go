package controllers

import (
	"go-booking-api/app/models"
	"go-booking-api/app/models/requests"
	"go-booking-api/app/services"
	"net/http"

	"github.com/revel/revel"
)

type AccountController struct {
	*revel.Controller
}

var us = services.GetUserService()
var vs = services.GetValidationService()
var rs = services.GetResponseService()

func (c AccountController) Login() revel.Result {

	var model requests.LoginRequest
	c.Params.BindJSON(&model)
	vs.ValidateLoginRequest(c.Controller, &model)

	if c.Validation.HasErrors() {
		c.Response.Status = http.StatusBadRequest
		resp := rs.CreateErrorResponse(c.Response.Status, "Login validation error", c.Validation.Errors)
		return c.RenderJSON(resp)
	}

	// go to db to check username and pwd
	// remember to hash pwd
	// if exists

	jwtService := &services.JwtService{}
	token := jwtService.GenerateToken(model.Email, true)
	resp := rs.CreateOperationResponse("jwt_token", token)
	return c.RenderJSON(resp)
}

func (c AccountController) Register() revel.Result {

	var model requests.RegisterRequest
	c.Params.BindJSON(&model)
	vs.ValidateRegisterRequest(c.Controller, &model)

	if c.Validation.HasErrors() {
		c.Response.Status = http.StatusBadRequest
		resp := rs.CreateErrorResponse(c.Response.Status, "Register validation error", c.Validation.Errors)
		return c.RenderJSON(resp)
	}

	if us.EmailExists(model.Email) {

		c.Response.Status = http.StatusBadRequest
		resp := rs.CreateErrorResponse(c.Response.Status, "Email exists", c.Validation.Errors)
		return c.RenderJSON(resp)
	}

	user := &models.User{ FirstName: model.FirstName, LastName: model.LastName, Email: model.Email,
		Credential: models.Credentials{Password: model.Password}}

	// save in db
	if err := us.Save(user); err != nil {

		c.Response.Status = http.StatusInternalServerError
		return c.Result
	}

	c.Response.Status = http.StatusCreated
	return c.Result
}
