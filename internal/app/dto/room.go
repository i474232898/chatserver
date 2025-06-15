package dto

import "time"

type CreateRoomRequest struct {
	Name string `json:"name" validate:"required"`
}

type RoomDTO struct {
	RoomId    uint
	RoomName  string
	CreatedAt *time.Time
	IsDirect  bool
	// Admin
	// Users
}

type CreateRoomDTO struct {
	AdminID   uint
	MemberIDs *[]int64
	CreateRoomRequest
}
