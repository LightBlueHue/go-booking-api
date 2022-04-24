package services

import (
	"go-booking-api/app/models"
)

type IUserService interface {
	EmailExists(email string) bool
	Save(user *models.User) error
	Get(email string) (*models.User, error)
	GetPassword(email string) (string, error)
}

type UserService struct {
}

func GetUserService() IUserService {

	return &UserService{}
}

func (service *UserService) EmailExists(email string) bool {

	db := GetDBService().GetDB()

	var age int
	db.Raw("SELECT COUNT(id) FROM users WHERE email = ?", email).Scan(&age)
	return age > 0
}

func (service *UserService) GetPassword(email string) (string, error) {

	db := GetDBService().GetDB()

	var pwd string
	result := db.Raw("SELECT password FROM users INNER JOIN credentials ON users.credential_id = credentials.id WHERE email = ?", email).Scan(&pwd)
	return pwd, result.Error
}

func (service *UserService) Save(user *models.User) error {

	db := GetDBService().GetDB()

	result := db.Create(user)
	return result.Error
}

func (service *UserService) Get(email string) (*models.User, error) {

	db := GetDBService().GetDB()

	var user models.User
	result := db.Where("email = ?", email).First(&user)
	return &user, result.Error
}
