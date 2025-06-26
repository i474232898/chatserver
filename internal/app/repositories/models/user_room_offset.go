package models

import (
	"time"
)

type UserRoomOffset struct {
	UserId          uint64 `gorm:"primaryKey;not null"`
	RoomId          uint64 `gorm:"primaryKey;not null"`
	LastReadMessage uint64 `gorm:"not null"`
	UpdatedAt       time.Time
}
