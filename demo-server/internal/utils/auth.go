package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserNumber string `json:"user_number"`
	Exp        int64  `json:"exp"`
	jwt.RegisteredClaims
}

func GenJwt(user_number string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_number": user_number,
		"exp":         time.Now().Add(time.Hour).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString(ENV["SERVER_SECRET"])
}

func CheckJwt(token string, userNumber string) error {

	parsedToken, err := jwt.ParseWithClaims(token, CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(ENV["SERVER_KEY"]), nil
	}, jwt.WithValidMethods([]string{"HS256"}))
	if err != nil {
		return err
	}

	if !parsedToken.Valid {
		return fmt.Errorf("invalid token")
	}

	claims, ok := parsedToken.Claims.(*CustomClaims)
	if !ok {
		return fmt.Errorf("invalid token")
	}

	if claims.UserNumber != userNumber {
		return fmt.Errorf("invalid token")
	}

	now := time.Now().Unix()
	if now > claims.Exp {
		return fmt.Errorf("invalid token")
	}

	return nil
}
