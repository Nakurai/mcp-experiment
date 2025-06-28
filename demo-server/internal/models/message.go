package models

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	ID           uint
	UserNumber   string
	Content      string
	ContactEmail string
	Tag          string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func CreateMessage(db *gorm.DB, msg *Message) (*Message, error) {
	result := db.Create(&msg)
	if result.Error != nil {
		return nil, result.Error
	}
	return msg, nil
}
