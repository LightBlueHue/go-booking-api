package services

import (
	"go-booking-api/app/models"
)

type IUserService interface {
	EmailExists(email string) bool
	Save(user *models.User) error
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

func (service *UserService) Save(user *models.User) error {

	db := GetDBService().GetDB()

	result := db.Create(user)
	return result.Error
}
