package controllers

import (
	"go-booking-api/app/models"
	"go-booking-api/app/models/requests"
	"go-booking-api/app/services"
	"net/http"
	"strings"

	"github.com/revel/revel"
)

type AccountController struct {
	*revel.Controller
}

func (c *AccountController) Login() revel.Result {

	_, hs, jwts, rs, us, vs := GetServices()
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
		c.Validation.Errors = append(c.Validation.Errors, &revel.ValidationError{Message: "hash", Key: "account"})
		resp := rs.CreateErrorResponse(c.Response.Status, "Login validation error", c.Validation.Errors)
		return c.RenderJSON(resp)
	}

	emailAndPwdExists, err := us.EmailAndPwdExists(model.Email, hashedPwd)

	if err != nil {

		c.Response.Status = http.StatusInternalServerError
		c.Validation.Errors = append(c.Validation.Errors, &revel.ValidationError{Message: "db", Key: "account"})
		resp := rs.CreateErrorResponse(c.Response.Status, "Sorry, We encountered an issue. Please try again", c.Validation.Errors)
		return c.RenderJSON(resp)
	}

	if !emailAndPwdExists {

		c.Response.Status = http.StatusBadRequest
		c.Validation.Errors = append(c.Validation.Errors, &revel.ValidationError{Message: "unknown account", Key: "account"})
		resp := rs.CreateErrorResponse(c.Response.Status, "Sorry, we can't find an account with this email address. Please try again or create a new account.", c.Validation.Errors)
		return c.RenderJSON(resp)
	}

	token := jwts.GenerateToken(model.Email, true)
	resp := rs.CreateOperationResponse("jwt_token", token)
	return c.RenderJSON(resp)
}

func (c *AccountController) Register() revel.Result {

	_, _, _, rs, us, vs := GetServices()
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

func IsLoggedIn(c revel.Controller) revel.Result {

	_, _, jwts, rs, _, _ := GetServices()
	auth := c.Request.Header.Get("Authorization")
	token := strings.Split(auth, " ")[1]
	jwtToken, err := jwts.ValidateToken(token)
	if err != nil || !jwtToken.Valid {

		c.Response.Status = http.StatusUnauthorized
		c.Validation.Errors = append(c.Validation.Errors, &revel.ValidationError{Message: "log in", Key: "account"})
		resp := rs.CreateErrorResponse(c.Response.Status, "Please log in", c.Validation.Errors)
		return c.RenderJSON(resp)
	}
	return nil
}

func GetServices() (services.IDBService, services.IHashService, services.IJWTService, services.IResponseService, services.IUserService, services.IValidationService) {

	return &services.DBService{}, &services.HashService{}, &services.JwtService{}, &services.ResponseService{}, &services.UserService{}, &services.ValidationService{}
}
