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

func NewService(db *gorm.DB) Service {

	return Service{DBService: NewDBService(db), HashService: NewHashService(), JWTService: NewJWTService(), ResponseService: NewResponseService(), UserService: NewUserService(db), ValidationService: GetValidationService(), BookingService: NewBookingService(db)}
}

// IsServiceInitialized returns true if all services are initialized.
func (s *Service) IsServiceInitialized() bool {

	return s.DBService != nil &&
		s.HashService != nil &&
		s.JWTService != nil &&
		s.ResponseService != nil &&
		s.UserService != nil &&
		s.ValidationService != nil &&
		s.BookingService != nil
}

type IDBService interface {
	InitDB(dbInfo DbInfo, dbInitializer DbInitializer, createDbStatement string) *gorm.DB
}

type (
	DbInitializer func(dialector gorm.Dialector, opts ...gorm.Option) (db *gorm.DB, err error)
)

type IHashService interface {
	HashAndSalt(password string) (string, error)
	CompareHashAndPassword(hashedPwd string, password string) (bool, error)
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
	GetByToken(token string, jwtService IJWTService) (*models.User, error)
}

type IValidationService interface {
	ValidateLoginRequest(v *revel.Validation, l *requests.LoginRequest)
	ValidateRegisterRequest(v *revel.Validation, l *requests.RegisterRequest)
	ValidateBookingRequest(v *revel.Validation, count uint)
}

type IBookingService interface {
	Book(userId uint, count uint) (uint, error)
	GetBookings(userId uint) (*[]models.Booking, error)
}
