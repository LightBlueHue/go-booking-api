package controllers

import (
	"go-booking-api/app/models/requests"
	"go-booking-api/app/services"
	"net/http"

	"github.com/revel/revel"
)

type AccountController struct {
	*revel.Controller
}

func (c AccountController) Login() revel.Result {

	var model requests.LoginRequest
	c.Params.BindJSON(&model)
	services.ValidateLoginRequest(c.Controller, &model)

	if c.Validation.HasErrors() {
		c.Response.Status = http.StatusBadRequest
		resp := services.CreateErrorResponse(c.Response.Status, "Login validation error", c.Validation.Errors)
		return c.RenderJSON(resp)
	}

	// go to db to check username and pwd
	// remember to hash pwd
	// if exists

	jwtService := &services.JwtService{}
	token := jwtService.GenerateToken(model.Email, true)
	resp := services.CreateOperationResponse("jwt_token", token)
	return c.RenderJSON(resp)
}

func (c AccountController) Register() revel.Result {

	var model requests.RegisterRequest
	c.Params.BindJSON(&model)
	services.ValidateRegisterRequest(c.Controller, &model)

	if c.Validation.HasErrors() {
		c.Response.Status = http.StatusBadRequest
		resp := services.CreateErrorResponse(c.Response.Status, "Register validation error", c.Validation.Errors)
		return c.RenderJSON(resp)
	}

	// save in db
	c.Response.Status = http.StatusCreated
	return c.Result
}
