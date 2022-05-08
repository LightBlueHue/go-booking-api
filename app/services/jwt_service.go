package services

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	EMAIL_CLAIM           JwtClaimType = "email"
	GO_BOOKING_API_SECRET string       = "GO_BOOKING_API_SECRET"
)

type JwtClaimType string

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

func GetJWTService() IJWTService {

	return &JwtService{
		secretKey: getSecretKey(),
		issure:    "go-booking-api",
	}
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

func (service *JwtService) GetClaim(token string, claimType JwtClaimType) (string, error) {

	jwtToken, err := service.ValidateToken(token)
	if err != nil || !jwtToken.Valid {

		return "", err
	}

	if claims, ok := jwtToken.Claims.(jwt.MapClaims); ok {

		if claim := claims[string(claimType)]; claim == nil {

			return "", fmt.Errorf("Unable to retrieve claim")
		}
		return fmt.Sprint(claims[string(claimType)]), nil
	}

	return "", fmt.Errorf("Unable to retrieve claim")
}

func getSecretKey() string {

	secret := os.Getenv(GO_BOOKING_API_SECRET)
	if secret == "" {

		panic(GO_BOOKING_API_SECRET)
	}

	return secret
}
