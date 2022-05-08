package services

import (
	"go-booking-api/app/models"
	"go-booking-api/app/models/requests"
	"go-booking-api/app/models/responses"

	"github.com/golang-jwt/jwt"
	"github.com/revel/revel"
	"gorm.io/gorm"
)

type Service struct {
	DBService         IDBService
	HashService       IHashService
	JWTService        IJWTService
	ResponseService   IResponseService
	UserService       IUserService
	ValidationService IValidationService
	BookingService    IBookingService
}

func (s *Service) SetServices(dBService IDBService, hashService IHashService, jwtService IJWTService, responseService IResponseService, userService IUserService, validationService IValidationService, bookingService IBookingService) {

	s.DBService = dBService
	s.HashService = hashService
	s.JWTService = jwtService
	s.ResponseService = responseService
	s.UserService = userService
	s.ValidationService = validationService
	s.BookingService = bookingService
}

type IService interface {
	SetServices(IDBService, IHashService, IJWTService, IResponseService, IUserService, IValidationService, IBookingService)
}

type IDBService interface {
	InitDB(database *gorm.DB, dbInfo DbInfo)
	GetDB() *gorm.DB
}

type IHashService interface {
	HashAndSalt(password string) (string, error)
	ComparePasswords(hashedPwd string, password string) bool
}

type IJWTService interface {
	GenerateToken(email string, isUser bool) string
	ValidateToken(token string) (*jwt.Token, error)
	GetClaim(token string, claimType JwtClaimType) (string, error)
}

type IResponseService interface {
	CreateErrorResponse(code int, message string, validationErrors []*revel.ValidationError) *responses.ErrorResponse
	CreateOperationResponse(context string, value interface{}) *responses.OperationResponse
}

type IUserService interface {
	EmailExists(email string) bool
	Save(user *models.User) error
	GetByEmail(email string) (*models.User, error)
	GetPassword(email string) (string, error)
	GetByToken(token string) (*models.User, error)
}

type IValidationService interface {
	ValidateLoginRequest(v *revel.Validation, l *requests.LoginRequest)
	ValidateRegisterRequest(v *revel.Validation, l *requests.RegisterRequest)
	ValidateBookingRequest(v *revel.Validation, count uint)
}