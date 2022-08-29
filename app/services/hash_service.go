package services

import (
	"golang.org/x/crypto/bcrypt"
)

type HashService struct{}

func NewHashService() IHashService {

	return &HashService{}
}

// HashAndSalt hashes and salts a password.
func (s *HashService) HashAndSalt(password string) (string, error) {

	pwd := []byte(password)
	// Use GenerateFromPassword to hash & salt pwd
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	if err != nil {

		return "", err
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash), nil
}

// CompareHashAndPassword compares a bcrypt hashed password with its possible plaintext equivalent. Returns true on success
func (s *HashService) CompareHashAndPassword(hashedPwd string, password string) (bool, error) {

	plainPwd := []byte(password)
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)

	if err := bcrypt.CompareHashAndPassword(byteHash, plainPwd); err != nil {

		return false, err
	}

	return true, nil
}
