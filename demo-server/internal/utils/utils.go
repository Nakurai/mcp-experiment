package utils

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"

	"gorm.io/gorm"
)

var ENV = map[string]string{}

func init() {
	// Load environement variables
	content, err := os.ReadFile(".env")
	if err != nil {
		log.Fatal(err)
	}
	content_str := string(content)
	envs := strings.Split(content_str, "\n")
	for _, env := range envs {
		parts := strings.Split(env, "=")
		ENV[parts[0]] = parts[1]
	}

}

type contextKey string

const DbSessionKey contextKey = "dbSession"

func GetDb(ctx context.Context) (*gorm.DB, error) {
	db, ok := ctx.Value(DbSessionKey).(*gorm.DB)
	if !ok {
		return nil, fmt.Errorf("no db session provided")
	}
	return db, nil
}

var allowedChar = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
var nbAllowedChar = big.NewInt(int64(len(allowedChar)))

// Get a random string composed of a-zA-Z0-9
func GetRandomString(length int) (string, error) {
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

func makeRes(w http.ResponseWriter, res any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
