package repositories

import (
	"context"

	// "github.com/i474232898/chatserver/internal/app/dto"
	"github.com/i474232898/chatserver/internal/app/repositories/models"
	"gorm.io/gorm"
)

type RoomRepository interface {
	Create(ctx context.Context, room *models.Room) (*models.Room, error)
	GetByName(ctx context.Context, name string) (*models.Room, error)
	GetById(ctx context.Context, id int64) (*models.Room, error)
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
