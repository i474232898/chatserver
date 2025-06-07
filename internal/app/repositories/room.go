package repositories

import (
	"context"

	"github.com/i474232898/chatserver/internal/app/repositories/models"
	"gorm.io/gorm"
)

type RoomRepository interface {
	Create(ctx context.Context, room *models.Room) (*models.Room, error)
	CreateDirectRoom(ctx context.Context, room *models.Room) (*models.Room, error)
	GetByName(ctx context.Context, name string) (*models.Room, error)
	GetById(ctx context.Context, id int64) (*models.Room, error)
	RoomsList(ctx context.Context, userId uint64) ([]models.Room, error)
	IsUserInRoom(ctx context.Context, userId uint64, roomId uint64) bool
}

type roomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) RoomRepository {
	return &roomRepository{db: db}
}

func (r *roomRepository) Create(ctx context.Context, room *models.Room) (*models.Room, error) {
	result := r.db.Create(room)

	if result.Error != nil {
		return nil, result.Error
	}

	return room, nil
}

func (r *roomRepository) CreateDirectRoom(ctx context.Context, room *models.Room) (*models.Room, error) {
	//check if room with same name exists
	var dbRoom models.Room
	if err := r.db.WithContext(ctx).Model(&models.Room{}).Where("name = ?", room.Name).First(&dbRoom).Error; err != nil {
		return nil, err
	}
	if dbRoom.ID != 0 {
		return &dbRoom, nil
	}

	return r.Create(ctx, room)
}

func (r *roomRepository) RoomsList(ctx context.Context, userId uint64) ([]models.Room, error) {
	var rooms []models.Room
	if err := r.db.WithContext(ctx).Where("id in (select room_id from rooms_users where user_id = ?)", userId).Find(&rooms).Error; err != nil {
		return nil, err
	}
	return rooms, nil
}

func (r *roomRepository) IsUserInRoom(ctx context.Context, userId uint64, roomId uint64) bool {
	var count int64
	db := r.db.WithContext(ctx).Model(&models.Room{})
	db = db.Where(
		"id = ? AND id in (select room_id from rooms_users where user_id = ?)",
		roomId, userId,
	)
	if err := db.Count(&count).Error; err != nil {
		return false
	}
	return count > 0
}

func (r *roomRepository) GetByName(ctx context.Context, name string) (*models.Room, error) {
	var room models.Room
	result := r.db.Where("name = ?", name).Preload("Admin").Preload("Users").First(&room)

	if result.Error != nil {
		return nil, result.Error
	}

	return &room, nil
}

func (r *roomRepository) GetById(ctx context.Context, id int64) (*models.Room, error) {
	var room models.Room
	result := r.db.Where("id = ?", id).Preload("Admin").Preload("Users").First(&room)

	if result.Error != nil {
		return nil, result.Error
	}

	return &room, nil
}
