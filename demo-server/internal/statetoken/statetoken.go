package statetoken

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

var tokenDb = map[string]int64{}

// Generate a new token you can use a unique code. Expires after 10 minutes.
// max length is 100
func GetNewToken(length int) (string, error) {
	if length > 100 {
		return "", fmt.Errorf("length too long. max length is 100")
	}
	randomString, err := getRandomString(length)
	if err != nil {
		return "", err
	}
	// the token expires in 10mn
	expiration := time.Now().Add(time.Minute * 10).Unix()
	tokenDb[randomString] = expiration

	return randomString, nil

}

// Check that the provided token is valid and not expired.
func VerifyToken(token string) error {
	expiration, ok := tokenDb[token]
	if !ok {
		return fmt.Errorf("token does not exist")
	}

	now := time.Now().Unix()
	if now > expiration {
		return fmt.Errorf("token is expired")
	}

	// here the token is valid so we remove it from the db
	delete(tokenDb, token)

	return nil
}

// Start a loop that will trigger every so often and remove expired tokens from db
func Init(ctx context.Context) {
	fmt.Println("starting state token loop")
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("closing state token loop")
				return
			case <-time.After(time.Minute):
				now := time.Now().Unix()
				for token, expiration := range tokenDb {
					if now > expiration {
						delete(tokenDb, token)
					}
				}
			}
		}
	}()
}

var allowedChar = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
var nbAllowedChar = big.NewInt(int64(len(allowedChar)))

// Get a random string composed of a-zA-Z0-9
func getRandomString(length int) (string, error) {
	randomString := ""
	for i := 0; i < length; i++ {
		charIndex, err := rand.Int(rand.Reader, nbAllowedChar)
		if err != nil {
			return "", fmt.Errorf("while generating random character %w", err)
		}
		newChar := string(allowedChar[charIndex.Int64()])
		randomString += newChar

	}

	return randomString, nil

}
