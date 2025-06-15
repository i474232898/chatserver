package dto

import "time"

// type CreateRoomRequest struct {
// 	Name string `json:"name" validate:"required,min=3,max=255"`
// }

type RoomDTO struct {
	RoomId    uint
	RoomName  string
	CreatedAt *time.Time
	IsDirect  bool
	// Admin
	// Users
}

type CreateRoomDTO struct {
	Name      string `json:"name" validate:"required,min=3,max=255"`
	AdminID   uint
	MemberIDs []int64 `json:"memberIDs" validate:"required,min=1,dive,gt=0"`
}

type CreateDirectRoomRequest struct {
	UserID int64 `json:"userID" validate:"required,gt=0"`
}
