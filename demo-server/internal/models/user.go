package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/nakurai/mcp-experiment/demo-server/internal/utils"
	"gorm.io/gorm"
)

type User struct {
	ID         uint
	Email      string `gorm:"uniqueIndex"`
	UserNumber string `gorm:"uniqueIndex"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func GetUserByEmail(db *gorm.DB, email string) (*User, error) {
	user := User{}
	result := db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return &user, nil

}

func CreateUser(db *gorm.DB, email string) (*User, error) {
	// saving the user in the database if they do not exist yet
	userNumber, err := utils.GetRandomString(20)
	if err != nil {
		return nil, fmt.Errorf("error generating new user number %w", err)
	}
	user := User{Email: email, UserNumber: "user-" + userNumber}
	result := db.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
