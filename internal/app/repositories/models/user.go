package models

import (
	// "time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"uniqueIndex;not null"`
	Username *string
	Password string
	// CreatedAt time.Time `gorm:"column:createdat;autoCreateTime"`
}
