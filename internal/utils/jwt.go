package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("123")

func GenerateToken(username string) (string, error) {
	fmt.Printf("username is %+v \n", username)
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	fmt.Printf("token is %+v \n", token)
	return token.SignedString(jwtSecret)
}
