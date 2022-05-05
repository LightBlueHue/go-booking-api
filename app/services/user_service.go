package services

import (
	"go-booking-api/app/models"
)

const (
	SQL_STATEMENT_GET_USER_PASSWORD = "SELECT password FROM users INNER JOIN credentials ON users.credential_id = credentials.id WHERE email = ?"
	SQL_STATEMENT_GET_EMAIL_COUNT   = "SELECT COUNT(id) FROM users WHERE email = ?"
	SQL_STATEMENT_GET_USER_BY_EMAIL = "email = ?"
)

type IUserService interface {
	EmailExists(email string) bool
	Save(user *models.User) error
	GetByEmail(email string) (*models.User, error)
	GetPassword(email string) (string, error)
	GetByToken(token string) (*models.User, error)
}

type UserService struct {
}

func GetUserService() IUserService {

	return &UserService{}
}

func (service *UserService) EmailExists(email string) bool {

	db := GetDBService().GetDB()

	var count int
	db.Raw(SQL_STATEMENT_GET_EMAIL_COUNT, email).Scan(&count)
	return count > 0
}

func (service *UserService) GetPassword(email string) (string, error) {

	db := GetDBService().GetDB()

	var pwd string
	result := db.Raw(SQL_STATEMENT_GET_USER_PASSWORD, email).Scan(&pwd)
	return pwd, result.Error
}

func (service *UserService) Save(user *models.User) error {

	db := GetDBService().GetDB()

	result := db.Create(user)
	return result.Error
}

func (service *UserService) GetByEmail(email string) (*models.User, error) {

	db := GetDBService().GetDB()

	var user models.User
	result := db.Where(SQL_STATEMENT_GET_USER_BY_EMAIL, email).First(&user)
	return &user, result.Error
}

func (service *UserService) GetByToken(token string) (*models.User, error) {

	var user *models.User
	var err error
	var email string
	if email, err = GetJWTService().GetClaim(token, EMAIL_CLAIM); err != nil {

		return nil, err
	}

	user, err = service.GetByEmail(email)
	return user, err
}
