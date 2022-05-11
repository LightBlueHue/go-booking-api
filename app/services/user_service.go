package services

import (
	"go-booking-api/app/models"

	"gorm.io/gorm"
)

const (
	SQL_STATEMENT_GET_USER_PASSWORD = "SELECT password FROM users INNER JOIN credentials ON users.credential_id = credentials.id WHERE email = ?"
	SQL_STATEMENT_GET_EMAIL_COUNT   = "SELECT COUNT(id) FROM users WHERE email = ?"
	SQL_STATEMENT_GET_USER_BY_EMAIL = "email = ?"
)

type UserService struct {
	db *gorm.DB
}

func GetUserService(db *gorm.DB) IUserService {

	return &UserService{db}
}

func (s *UserService) EmailExists(email string) bool {

	var count int
	s.db.Raw(SQL_STATEMENT_GET_EMAIL_COUNT, email).Scan(&count)
	return count > 0
}

func (s *UserService) GetPassword(email string) (string, error) {

	var pwd string
	result := s.db.Raw(SQL_STATEMENT_GET_USER_PASSWORD, email).Scan(&pwd)
	return pwd, result.Error
}

func (s *UserService) Save(user *models.User) error {

	result := s.db.Create(user)
	return result.Error
}

func (s *UserService) GetByEmail(email string) (*models.User, error) {

	var user models.User
	result := s.db.Where(SQL_STATEMENT_GET_USER_BY_EMAIL, email).First(&user)

	if result.Error != nil {

		return nil, result.Error
	}

	return &user, result.Error
}

func (s *UserService) GetByToken(token string, jwtService IJWTService) (*models.User, error) {

	var user *models.User
	var err error
	var email string
	if email, err = jwtService.GetClaim(token, EMAIL_CLAIM); err != nil {

		return nil, err
	}

	user, err = s.GetByEmail(email)
	return user, err
}
