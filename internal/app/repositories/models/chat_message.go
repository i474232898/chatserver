package models

import "gorm.io/gorm"

type ChatMessage struct {
	gorm.Model
	RoomId  uint64 `gorm:"not null"`
	UserId  uint64 `gorm:"not null"`
	Content string `gorm:"not null"`
}
