package dto

import "time"

type CreateRoomRequest struct {
	Name string `json:"name" validate:"required"`
}

type RoomDTO struct {
	ID        uint
	Name      string
	CreatedAt *time.Time
	// Admin
	// Users
}

type NewRoomDTO struct {
	AdminID uint
	CreateRoomRequest
}
