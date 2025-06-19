package repositories

import (
	"context"
	"fmt"

	"github.com/i474232898/chatserver/internal/app/repositories/models"
	"gorm.io/gorm"
)

type messageRepository struct {
	db *gorm.DB
}

type MessageRepository interface {
	Create(ctx context.Context, msg *models.ChatMessage) (models.ChatMessage, error)
	GetMessages(ctx context.Context, roomId, lastMessageId uint64) ([]models.ChatMessage, error)
}

func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &messageRepository{db: db}
}

func (r messageRepository) Create(ctx context.Context, msg *models.ChatMessage) (models.ChatMessage, error) {
	result := r.db.WithContext(ctx).Create(msg)
	if result.Error != nil {
		return models.ChatMessage{}, result.Error
	}
	return *msg, nil
}

func (r messageRepository) GetMessages(ctx context.Context, roomId, lastMessageId uint64) ([]models.ChatMessage, error) {
	var msgs []models.ChatMessage
	fmt.Println(lastMessageId, "<<lastMessageId")

	if lastMessageId == 0 {
		result := r.db.WithContext(ctx).Where("room_id = ?", roomId).Find(&msgs)
		if result.Error != nil {
			return nil, result.Error
		}
		return msgs, nil
	}

	createdAt := r.db.Model(&models.ChatMessage{}).
		Select("COALESCE(created_at, '1970-01-01 00:00:00')").Where("id=?", lastMessageId)
	mainQuery := r.db.WithContext(ctx).Where("room_id = ? AND created_at > (?)", roomId, createdAt)

	sql := mainQuery.ToSQL(func(tx *gorm.DB) *gorm.DB { return tx.Find(&msgs) })
	fmt.Println(sql)

	result := mainQuery.Find(&msgs)
	if result.Error != nil {
		return nil, result.Error
	}

	return msgs, nil
}
