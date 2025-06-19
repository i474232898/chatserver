package repositories

import (
	"context"

	"github.com/i474232898/chatserver/internal/app/repositories/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type userRoomOffsetRepository struct {
	db *gorm.DB
}

type UserRoomOffsetRepository interface {
	UpdateUserRoomOffset(ctx context.Context, roomId uint64, userId uint64, lastReadMessage uint64) error
	GetUserRoomOffset(ctx context.Context, roomId uint64, userId uint64) (uint64, error)
}

func NewUserRoomOffsetRepository(db *gorm.DB) UserRoomOffsetRepository {
	return &userRoomOffsetRepository{db: db}
}

func (r *userRoomOffsetRepository) UpdateUserRoomOffset(ctx context.Context, roomId uint64,
	userId uint64, lastReadMessage uint64) error {

	result := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "room_id"}, {Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"last_read_message"}),
	}).Create(&models.UserRoomOffset{
		RoomId:          roomId,
		UserId:          userId,
		LastReadMessage: lastReadMessage,
	})

	return result.Error
}

func (r *userRoomOffsetRepository) GetUserRoomOffset(ctx context.Context,
	roomId, userId uint64) (uint64, error) {
	var lastReadMessage uint64
	result := r.db.WithContext(ctx).Model(&models.UserRoomOffset{}).
		Where("room_id = ? AND user_id = ?", roomId, userId).
		Select("last_read_message").First(&lastReadMessage)
	if result.Error != nil {
		return 0, result.Error
	}
	return lastReadMessage, nil
}
