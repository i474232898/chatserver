package validations

import (
	"github.com/i474232898/chatserver/api/types"
	"github.com/i474232898/chatserver/internal/app/dto"
)

func ValidateCreateRoom(room *types.CreateRoomRequest) error {
	data := dto.CreateRoomDTO{
		Name:      room.Name,
		MemberIDs: room.MemberIDs,
	}

	return Validator.Struct(data)
}

func ValidateCreateDirectRoom(room *types.CreateDirectRoomRequest) error {
	data := dto.CreateDirectRoomRequest{
		UserID: room.UserID,
	}

	return Validator.Struct(data)
}
