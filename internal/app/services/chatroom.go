package services

import (
	"context"
	"sort"
	"strconv"

	// "github.com/i474232898/chatserver/api/types"
	"github.com/i474232898/chatserver/api/types"
	"github.com/i474232898/chatserver/internal/app/dto"
	"github.com/i474232898/chatserver/internal/app/repositories"
	"github.com/i474232898/chatserver/internal/app/repositories/models"
	"gorm.io/gorm"
)

type ChatRoomService interface {
	Create(ctx context.Context, room *dto.NewRoomDTO) (*dto.RoomDTO, error)
	GetByName(ctx context.Context, name string) (*models.Room, error)
	GetById(ctx context.Context, id int64) (*models.Room, error)
	ListRooms(ctx context.Context, userId uint64) (types.RoomsListResponse, error)
}

type chatRoomService struct {
	roomRepo repositories.RoomRepository
}

func NewChatRoomService(roomRepo repositories.RoomRepository) ChatRoomService {
	return &chatRoomService{roomRepo: roomRepo}
}

func (s *chatRoomService) Create(ctx context.Context, room *dto.NewRoomDTO) (*dto.RoomDTO, error) {
	users := []models.User{}
	for _, id := range *room.MemberIDs {
		users = append(users, models.User{Model: gorm.Model{ID: uint(id)}})
	}
	users = append(users, models.User{Model: gorm.Model{ID: room.AdminID}})

	roomName := room.Name
	//direct message
	if len(*room.MemberIDs) == 1 && room.Name == "" {
		roomName = generateDirectRoomName(*room.MemberIDs)
	}

	newRoom := &models.Room{
		Name:    roomName,
		AdminID: room.AdminID,
		Users:   users,
	}

	_, err := s.roomRepo.Create(ctx, newRoom)
	if err != nil {
		return nil, err
	}

	return &dto.RoomDTO{
		RoomName:  newRoom.Name,
		RoomId:    newRoom.ID,
		CreatedAt: &newRoom.CreatedAt,
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

func (s *chatRoomService) GetByName(ctx context.Context, name string) (*models.Room, error) {
	return s.roomRepo.GetByName(ctx, name)
}
func (s *chatRoomService) GetById(ctx context.Context, id int64) (*models.Room, error) {
	return s.roomRepo.GetById(ctx, id)
}

func generateDirectRoomName(userIDs []int64) string {
	sort.Slice(userIDs, func(i, j int) bool {
		return userIDs[i] < userIDs[j]
	})

	return "direct-" + strconv.Itoa(int(userIDs[0])) + "-" + strconv.Itoa(int(userIDs[1]))
}
