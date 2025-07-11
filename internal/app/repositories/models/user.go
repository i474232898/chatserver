package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"uniqueIndex;not null"`
	Username *string
	Password string
	Rooms    []Room `gorm:"many2many:rooms_users;"`
}
