package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/nakurai/mcp-experiment/demo-server/internal/utils"
	"gorm.io/gorm"
)

type AuthToken struct {
	ID         uint
	UserNumber string
	Token      string `gorm:"uniqueIndex"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

// Get the requested token
func GetToken(db *gorm.DB, token string) (*AuthToken, error) {
	authToken := AuthToken{}
	result := db.Where("token = ?", token).First(&authToken)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return &authToken, nil

}

// given a user number, generate a new JWT for them and store it in the database
func CreateToken(db *gorm.DB, userNumber string) (*AuthToken, error) {
	token, err := utils.GenJwt(userNumber)
	if err != nil {
		return nil, fmt.Errorf("error generating new jwt %w", err)
	}
	authToken := AuthToken{Token: token, UserNumber: userNumber}
	result := db.Create(&authToken)
	if result.Error != nil {
		return nil, result.Error
	}
	return &authToken, nil
}
