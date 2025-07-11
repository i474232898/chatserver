package models

import (
	"gorm.io/gorm"
)

type Room struct {
	gorm.Model
	Name     string
	AdminID  uint
	Admin    User
	Users    []User `gorm:"many2many:rooms_users;"`
	IsDirect bool
}
