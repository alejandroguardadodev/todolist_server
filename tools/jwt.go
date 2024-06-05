package tools

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWTToken(userId uint, userType string, contact string, contactType string) (string, error) {
	secretKey := os.Getenv("SECRET_KEY")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":      userId,
		"user_type":    userType,
		"contact":      contact,
		"contact_type": contactType,
		"exp":          time.Now().Add(time.Hour * 2).Unix(), // The token only gonna last 2 hours
	})

	return token.SignedString([]byte(secretKey))
}

func VerifyJWTToken(token string) (jwt.MapClaims, error) {
	secretKey := os.Getenv("SECRET_KEY")

	pasedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errors.New("unexpected signed method")
		}

		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, errors.New("token could not be parsed")
	}

	tokenIsValid := pasedToken.Valid

	if !tokenIsValid {
		return nil, errors.New("invalid token")
	}

	claims, ok := pasedToken.Claims.(jwt.MapClaims)

	if !ok {
		return nil, errors.New("claims token error")
	}

	return claims, nil
}
