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
var jwts = services.GetJWTService()
var hs = services.GetHashService()

func (c AccountController) Login() revel.Result {

	var model requests.LoginRequest
	c.Params.BindJSON(&model)
	vs.ValidateLoginRequest(c.Controller, &model)

	if c.Validation.HasErrors() {

		c.Response.Status = http.StatusBadRequest
		resp := rs.CreateErrorResponse(c.Response.Status, "Login validation error", c.Validation.Errors)
		return c.RenderJSON(resp)
	}

	var hashedPwd string
	var err error
	if hashedPwd, err = hs.HashAndSalt(model.Password); err != nil {

		c.Response.Status = http.StatusBadRequest
		c.Validation.Errors = append(c.Validation.Errors, &revel.ValidationError{Message: "hash", Key: "Account"})
		resp := rs.CreateErrorResponse(c.Response.Status, "Login validation error", c.Validation.Errors)
		return c.RenderJSON(resp)
	}

	emailAndPwdExists, err := us.EmailAndPwdExists(model.Email, hashedPwd)

	if err != nil {

		c.Response.Status = http.StatusInternalServerError
		c.Validation.Errors = append(c.Validation.Errors, &revel.ValidationError{Message: "db", Key: "Account"})
		resp := rs.CreateErrorResponse(c.Response.Status, "Sorry, We encountered an issue. Please try again", c.Validation.Errors)
		return c.RenderJSON(resp)
	}

	if !emailAndPwdExists {

		c.Response.Status = http.StatusBadRequest
		c.Validation.Errors = append(c.Validation.Errors, &revel.ValidationError{Message: "unknown account", Key: "Account"})
		resp := rs.CreateErrorResponse(c.Response.Status, "Sorry, we can't find an account with this email address. Please try again or create a new account.", c.Validation.Errors)
		return c.RenderJSON(resp)
	}

	token := jwts.GenerateToken(model.Email, true)
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

	user := &models.User{FirstName: model.FirstName, LastName: model.LastName, Email: model.Email,
		Credential: models.Credentials{Password: model.Password}}

	// save in db
	if err := us.Save(user); err != nil {

		c.Response.Status = http.StatusInternalServerError
		return c.Result
	}

	c.Response.Status = http.StatusCreated
	return c.Result
}
