package services

import (
	"context"

	// "github.com/i474232898/chatserver/api/types"
	"github.com/i474232898/chatserver/internal/app/dto"
	"github.com/i474232898/chatserver/internal/app/repositories"
	"github.com/i474232898/chatserver/internal/app/repositories/models"
)

type ChatRoomService interface {
	Create(ctx context.Context, room *dto.NewRoomDTO) (*dto.RoomDTO, error)
	GetByName(ctx context.Context, name string) (*models.Room, error)
	GetById(ctx context.Context, id int64) (*models.Room, error)
}

type chatRoomService struct {
	roomRepo repositories.RoomRepository
}

func NewChatRoomService(roomRepo repositories.RoomRepository) ChatRoomService {
	return &chatRoomService{roomRepo: roomRepo}
}

func (s *chatRoomService) Create(ctx context.Context, room *dto.NewRoomDTO) (*dto.RoomDTO, error) {
	newRoom, err := s.roomRepo.Create(ctx, room)
	if err != nil {
		return nil, err
	}

	return &dto.RoomDTO{
		Name:      newRoom.Name,
		ID:        newRoom.ID,
		CreatedAt: &newRoom.CreatedAt,
	}, nil
}
func (s *chatRoomService) GetByName(ctx context.Context, name string) (*models.Room, error) {
	return s.roomRepo.GetByName(ctx, name)
}
func (s *chatRoomService) GetById(ctx context.Context, id int64) (*models.Room, error) {
	return s.roomRepo.GetById(ctx, id)
}
