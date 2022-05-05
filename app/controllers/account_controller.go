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
}

func (c *AccountController) Login() revel.Result {

	_, hs, jwts, rs, us, vs, _ := GetServices()
	var model requests.LoginRequest
	c.Params.BindJSON(&model)
	vs.ValidateLoginRequest(c.Controller, &model)

	if c.Validation.HasErrors() {

		c.Response.Status = http.StatusBadRequest
		resp := rs.CreateErrorResponse(c.Response.Status, "Login validation error", c.Validation.Errors)
		return c.RenderJSON(resp)
	}

	hashedPwd, err := us.GetPassword(model.Email)

	if err != nil {

		c.Response.Status = http.StatusInternalServerError
		c.Validation.Errors = append(c.Validation.Errors, &revel.ValidationError{Message: "db", Key: "account"})
		resp := rs.CreateErrorResponse(c.Response.Status, "Sorry, We encountered an issue. Please try again", c.Validation.Errors)
		return c.RenderJSON(resp)
	}

	if pwdEqual := hs.ComparePasswords(hashedPwd, model.Password); !pwdEqual {

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

	_, hs, _, rs, us, vs, _ := GetServices()
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

	var err error
	var hshPwd string
	if hshPwd, err = hs.HashAndSalt(model.Password); err != nil {

		c.Response.Status = http.StatusInternalServerError
		c.Validation.Errors = append(c.Validation.Errors, &revel.ValidationError{Message: "h&s", Key: "account"})
		resp := rs.CreateErrorResponse(c.Response.Status, "Please try again", c.Validation.Errors)
		return c.RenderJSON(resp)
	}

	user := &models.User{FirstName: model.FirstName, LastName: model.LastName, Email: model.Email,
		Credential: models.Credentials{Password: hshPwd}}

	// save in db
	if err = us.Save(user); err != nil {

		c.Response.Status = http.StatusInternalServerError
		c.Validation.Errors = append(c.Validation.Errors, &revel.ValidationError{Message: "save", Key: "account"})
		resp := rs.CreateErrorResponse(c.Response.Status, "Please try again", c.Validation.Errors)
		return c.RenderJSON(resp)
	}

	c.Response.Status = http.StatusCreated
	return c.Result
}

func IsLoggedIn(c *revel.Controller) revel.Result {

	_, _, jwts, rs, us, _, _ := GetServices()
	var token string
	var err error

	if token, err = getBearerToken(c); err != nil {

		c.Response.Status = http.StatusUnauthorized
		c.Validation.Errors = append(c.Validation.Errors, &revel.ValidationError{Message: err.Error(), Key: "account"})
		resp := rs.CreateErrorResponse(c.Response.Status, "Please log in", c.Validation.Errors)
		return c.RenderJSON(resp)
	}

	if jwtToken, err := jwts.ValidateToken(token); err != nil || !jwtToken.Valid {

		c.Response.Status = http.StatusUnauthorized
		c.Validation.Errors = append(c.Validation.Errors, &revel.ValidationError{Message: "log in", Key: "account"})
		resp := rs.CreateErrorResponse(c.Response.Status, "Please log in", c.Validation.Errors)
		return c.RenderJSON(resp)
	}

	var user *models.User
	if user, err = us.GetByToken(token); err != nil || user == nil {

		c.Validation.Errors = append(c.Validation.Errors, &revel.ValidationError{Message: "log in", Key: "account"})
		c.Response.Status = http.StatusInternalServerError
		resp := rs.CreateErrorResponse(c.Response.Status, "Please log in", c.Validation.Errors)
		return c.RenderJSON(resp)
	}

	return nil
}

func GetUser(c *revel.Controller) (*models.User, error) {

	_, _, jwts, _, us, _, _ := GetServices()
	var token string
	var err error

	if token, err = getBearerToken(c); err != nil {

		return nil, err
	}

	if jwtToken, err := jwts.ValidateToken(token); err != nil || !jwtToken.Valid {

		return nil, err
	}

	var user *models.User
	if user, err = us.GetByToken(token); err != nil || user == nil {

		if err == nil {

			return nil, fmt.Errorf("")
		}

		return nil, err
	}

	return user, nil
}

func GetServices() (services.IDBService, services.IHashService, services.IJWTService, services.IResponseService, services.IUserService, services.IValidationService, services.IBookingService) {

	return services.GetDBService(), services.GetHashService(), services.GetJWTService(), services.GetResponseService(), services.GetUserService(), services.GetValidationService(), services.GetBookingService()
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