package services

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

type IJWTService interface {
	GenerateToken(email string, isUser bool) string
	ValidateToken(token string) (*jwt.Token, error)
	GetClaim(token string, claim string) (string, error)
}

type authCustomClaims struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	User  bool   `json:"user"`
	jwt.StandardClaims
}

type JwtService struct {
	secretKey string
	issure    string
}

const (
	EmailClaimType = "email"
)

func GetJWTService() IJWTService {

	return &JwtService{
		secretKey: getSecretKey(),
		issure:    "go-booikng-api",
	}
}

func getSecretKey() string {

	secret := os.Getenv("SECRET")
	if secret == "" {

		secret = "secret"
	}
	return secret
}

func (service *JwtService) GenerateToken(email string, isUser bool) string {

	claims := &authCustomClaims{
		email,
		email,
		isUser,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 48).Unix(),
			Issuer:    service.issure,
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//encoded string
	t, err := token.SignedString([]byte(service.secretKey))
	if err != nil {

		panic(err)
	}
	return t
}

func (service *JwtService) ValidateToken(encodedToken string) (*jwt.Token, error) {

	return jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if _, isvalid := token.Method.(*jwt.SigningMethodHMAC); !isvalid {

			return nil, fmt.Errorf("Invalid token %v", token.Header["alg"])
		}
		return []byte(service.secretKey), nil
	})
}

func (service *JwtService) GetClaim(token string, claim string) (string, error) {

	jwtToken, err := service.ValidateToken(token)
	if err != nil || !jwtToken.Valid {

		return "", err
	}

	if claims, ok := jwtToken.Claims.(jwt.MapClaims); ok {

		return fmt.Sprint(claims[claim]), nil
	}

	return "", fmt.Errorf("Unable to retrieve claim")
}
