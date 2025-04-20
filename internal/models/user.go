package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique"`
	Email    string `gorm:"unique"`
	Uuid     string `gorm:"unique"`
}

type ChatMessage struct {
	SenderID    uint      `bson:"sender_id"`
	RecipientID uint      `bson:"recipient_id"`
	Content     string    `bson:"content"`
	Timestamp   time.Time `bson:"timestamp"`
}
