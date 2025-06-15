package services

import (
	"context"
	"sort"
	"strconv"

	"github.com/i474232898/chatserver/api/types"
	"github.com/i474232898/chatserver/internal/app/dto"
	"github.com/i474232898/chatserver/internal/app/repositories"
	"github.com/i474232898/chatserver/internal/app/repositories/models"
	"gorm.io/gorm"
)

type ChatRoomService interface {
	Create(ctx context.Context, room *dto.CreateRoomDTO) (*dto.RoomDTO, error)
	GetByName(ctx context.Context, name string) (*dto.RoomDTO, error)
	ListRooms(ctx context.Context, userId uint64) (types.RoomsListResponse, error)
	IsUserInRoom(ctx context.Context, userId uint64, roomId uint64) bool
}

type chatRoomService struct {
	roomRepo repositories.RoomRepository
}

func NewChatRoomService(roomRepo repositories.RoomRepository) ChatRoomService {
	return &chatRoomService{roomRepo: roomRepo}
}

func (s *chatRoomService) Create(ctx context.Context, room *dto.CreateRoomDTO) (*dto.RoomDTO, error) {
	users := []models.User{}
	for _, id := range *room.MemberIDs {
		users = append(users, models.User{Model: gorm.Model{ID: uint(id)}})
	}
	users = append(users, models.User{Model: gorm.Model{ID: room.AdminID}})

	roomName := room.Name
	//direct message
	if len(users) == 2 && room.Name == "" {
		roomName = generateDirectRoomName(users)
	}

	newRoom := &models.Room{
		Name:     roomName,
		AdminID:  room.AdminID,
		Users:    users,
		IsDirect: len(users) == 2,
	}

	var err error
	var roomDb *models.Room
	if len(users) == 2 {
		roomDb, err = s.roomRepo.CreateDirectRoom(ctx, newRoom)
	} else {
		roomDb, err = s.roomRepo.Create(ctx, newRoom)
	}
	if err != nil {
		return nil, err
	}

	return &dto.RoomDTO{
		RoomName:  roomDb.Name,
		RoomId:    roomDb.ID,
		CreatedAt: &newRoom.CreatedAt,
		IsDirect:  roomDb.IsDirect,
	}, nil
}

func (s *chatRoomService) ListRooms(ctx context.Context, userId uint64) (types.RoomsListResponse, error) {
	rooms, err := s.roomRepo.RoomsList(ctx, userId)
	if err != nil {
		return nil, err
	}
	roomsResp := types.RoomsListResponse{}

	for _, room := range rooms {
		roomsResp = append(roomsResp, types.Room{
			CreatedAt: room.CreatedAt,
			Id:        int64(room.ID),
			Name:      room.Name,
		})
	}
	return roomsResp, nil
}

func generateDirectRoomName(users []models.User) string {
	sort.Slice(users, func(i, j int) bool {
		return users[i].ID < users[j].ID
	})

	return "direct-" + strconv.Itoa(int(users[0].ID)) + "-" + strconv.Itoa(int(users[1].ID))
}

func (s *chatRoomService) IsUserInRoom(ctx context.Context, userId uint64, roomId uint64) bool {
	return s.roomRepo.IsUserInRoom(ctx, userId, roomId)
}

func (s *chatRoomService) GetByName(ctx context.Context, name string) (*dto.RoomDTO, error) {
	model, err := s.roomRepo.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}
	return &dto.RoomDTO{
		RoomName:  model.Name,
		RoomId:    model.ID,
		CreatedAt: &model.CreatedAt,
	}, nil
}
