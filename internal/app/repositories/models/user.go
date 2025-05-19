package models

import "time"

type User struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	Email     string    `gorm:"uniqueIndex;not null" json:"email"`
	Username  *string   `json:"username"`
	Password  string    `json:"password"`
	CreatedAt time.Time `gorm:"column:createdat;autoCreateTime" json:"createdAt"`
}
