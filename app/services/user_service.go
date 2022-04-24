package services

import (
	"go-booking-api/app/models"
)

type IUserService interface {
	EmailExists(email string) bool
	Save(user *models.User) error
	GetByEmailAndPwd(email string, hashedPwd string) (*models.User, error)
	EmailAndPwdExists(email string, hashedPwd string) (bool, error)
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

func (service *UserService) EmailAndPwdExists(email string, hashedPwd string) (bool, error) {

	db := GetDBService().GetDB()

	var count int
	result := db.Raw("SELECT COUNT(users.id) FROM users INNER JOIN credentials ON users.credential_id = credentials.id WHERE email = ? AND password =?", email, hashedPwd).Scan(&count)
	return count == 1, result.Error
}

func (service *UserService) Save(user *models.User) error {

	db := GetDBService().GetDB()

	result := db.Create(user)
	return result.Error
}

func (service *UserService) GetByEmailAndPwd(email string, hashedPwd string) (*models.User, error) {

	db := GetDBService().GetDB()

	var user models.User
	result := db.Where("email = ? AND password =?", email, hashedPwd).First(&user)
	return &user, result.Error
}
