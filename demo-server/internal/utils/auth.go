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
	exp := time.Now().Add(time.Hour).Unix()
	claims := CustomClaims{
		UserNumber: user_number,
		Exp:        exp,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString([]byte(ENV["SERVER_KEY"]))
}

func CheckJwt(token string, userNumber string) error {
	claims := CustomClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(ENV["SERVER_KEY"]), nil
	})
	if err != nil {
		return err
	}

	if !parsedToken.Valid {
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
