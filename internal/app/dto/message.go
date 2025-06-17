package dto

import "time"

type MessageDTO struct {
	ID        uint
	RoomId    uint64
	UserId    uint64
	Content   string
	CreatedAt time.Time
}
