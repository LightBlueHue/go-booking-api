package controllers

import (
	"fmt"
	"go-booking-api/app/models"
	"go-booking-api/app/models/requests"
	"go-booking-api/app/services"
	"net/http"
	"strings"

	"github.com/revel/revel"
)

type AccountController struct {
	*revel.Controller
	Service services.Service
}

func (c *AccountController) Login() revel.Result {

	var model requests.LoginRequest
	c.Params.BindJSON(&model)
	c.Service.ValidationService.ValidateLoginRequest(c.Validation, &model)

	if c.Validation.HasErrors() {

		c.Response.Status = http.StatusBadRequest
		resp := c.Service.ResponseService.CreateErrorResponse(c.Response.Status, "Login validation error", c.Validation.Errors)
		return c.RenderJSON(resp)
	}

	hashedPwd, err := c.Service.UserService.GetPassword(model.Email)

	if err != nil {

		c.Response.Status = http.StatusInternalServerError
		c.Validation.Errors = append(c.Validation.Errors, &revel.ValidationError{Message: "db", Key: "account"})
		resp := c.Service.ResponseService.CreateErrorResponse(c.Response.Status, "Sorry, We encountered an issue. Please try again", c.Validation.Errors)
		return c.RenderJSON(resp)
	}

	if pwdEqual := c.Service.HashService.ComparePasswords(hashedPwd, model.Password); !pwdEqual {

		c.Response.Status = http.StatusBadRequest
		c.Validation.Errors = append(c.Validation.Errors, &revel.ValidationError{Message: "unknown account", Key: "account"})
		resp := c.Service.ResponseService.CreateErrorResponse(c.Response.Status, "Sorry, we can't find an account with this email address. Please try again or create a new account.", c.Validation.Errors)
		return c.RenderJSON(resp)
	}

	token := c.Service.JWTService.GenerateToken(model.Email, true)
	resp := c.Service.ResponseService.CreateOperationResponse("jwt_token", token)
	return c.RenderJSON(resp)
}

func (c *AccountController) Register() revel.Result {

	var model requests.RegisterRequest
	c.Params.BindJSON(&model)
	c.Service.ValidationService.ValidateRegisterRequest(c.Validation, &model)

	if c.Validation.HasErrors() {

		c.Response.Status = http.StatusBadRequest
		resp := c.Service.ResponseService.CreateErrorResponse(c.Response.Status, "Register validation error", c.Validation.Errors)
		return c.RenderJSON(resp)
	}

	if c.Service.UserService.EmailExists(model.Email) {

		c.Response.Status = http.StatusBadRequest
		resp := c.Service.ResponseService.CreateErrorResponse(c.Response.Status, "Email exists", c.Validation.Errors)
		return c.RenderJSON(resp)
	}

	var err error
	var hshPwd string
	if hshPwd, err = c.Service.HashService.HashAndSalt(model.Password); err != nil {

		c.Response.Status = http.StatusInternalServerError
		c.Validation.Errors = append(c.Validation.Errors, &revel.ValidationError{Message: "h&s", Key: "account"})
		resp := c.Service.ResponseService.CreateErrorResponse(c.Response.Status, "Please try again", c.Validation.Errors)
		return c.RenderJSON(resp)
	}

	user := &models.User{FirstName: model.FirstName, LastName: model.LastName, Email: model.Email,
		Credential: models.Credentials{Password: hshPwd}}

	// save in db
	if err = c.Service.UserService.Save(user); err != nil {

		c.Response.Status = http.StatusInternalServerError
		c.Validation.Errors = append(c.Validation.Errors, &revel.ValidationError{Message: "save", Key: "account"})
		resp := c.Service.ResponseService.CreateErrorResponse(c.Response.Status, "Please try again", c.Validation.Errors)
		return c.RenderJSON(resp)
	}

	c.Response.Status = http.StatusCreated
	return c.Result
}

func IsLoggedIn(c *revel.Controller, s services.Service) revel.Result {

	var token string
	var err error

	if token, err = getBearerToken(c); err != nil {

		c.Response.Status = http.StatusUnauthorized
		c.Validation.Errors = append(c.Validation.Errors, &revel.ValidationError{Message: err.Error(), Key: "account"})
		resp := s.ResponseService.CreateErrorResponse(c.Response.Status, "Please log in", c.Validation.Errors)
		return c.RenderJSON(resp)
	}

	if jwtToken, err := s.JWTService.ValidateToken(token); err != nil || !jwtToken.Valid {

		c.Response.Status = http.StatusUnauthorized
		c.Validation.Errors = append(c.Validation.Errors, &revel.ValidationError{Message: "log in", Key: "account"})
		resp := s.ResponseService.CreateErrorResponse(c.Response.Status, "Please log in", c.Validation.Errors)
		return c.RenderJSON(resp)
	}

	var user *models.User
	if user, err = s.UserService.GetByToken(token); err != nil || user == nil {

		c.Validation.Errors = append(c.Validation.Errors, &revel.ValidationError{Message: "log in", Key: "account"})
		c.Response.Status = http.StatusInternalServerError
		resp := s.ResponseService.CreateErrorResponse(c.Response.Status, "Please log in", c.Validation.Errors)
		return c.RenderJSON(resp)
	}

	return nil
}

func GetUser(c *revel.Controller, s services.Service) (*models.User, error) {

	var token string
	var err error

	if token, err = getBearerToken(c); err != nil {

		return nil, err
	}

	if jwtToken, err := s.JWTService.ValidateToken(token); err != nil || !jwtToken.Valid {

		return nil, err
	}

	var user *models.User
	if user, err = s.UserService.GetByToken(token); err != nil || user == nil {

		if err == nil {

			return nil, fmt.Errorf("")
		}

		return nil, err
	}

	return user, nil
}

func getBearerToken(c *revel.Controller) (string, error) {

	auth := c.Request.Header.Get("Authorization")
	if strings.TrimSpace(auth) == "" {

		return "", fmt.Errorf("Authorization header empty")
	}

	var token string
	if token = strings.Split(auth, " ")[1]; token == "" {
		return "", fmt.Errorf("Authorization header Bearer empty")
	}

	return token, nil
}
