package services

import (
	"golang.org/x/crypto/bcrypt"
)

type HashService struct{}

func GetHashService() IHashService {

	return &HashService{}
}

func (s *HashService) HashAndSalt(password string) (string, error) {

	pwd := []byte(password)
	// Use GenerateFromPassword to hash & salt pwd
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {

		return "", err
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash), nil
}

func (s *HashService) ComparePasswords(hashedPwd string, password string) bool {

	plainPwd := []byte(password)
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)

	if err := bcrypt.CompareHashAndPassword(byteHash, plainPwd); err != nil {

		return false
	}

	return true
}
