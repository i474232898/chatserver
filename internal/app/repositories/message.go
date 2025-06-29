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
	Create(ctx context.Context, msg models.ChatMessage) (models.ChatMessage, error)
	GetMessages(ctx context.Context, roomId, lastMessageId uint64) ([]models.ChatMessage, error)
}

func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &messageRepository{db: db}
}

func (r messageRepository) Create(ctx context.Context, msg models.ChatMessage) (models.ChatMessage, error) {
	result := r.db.WithContext(ctx).Create(msg)
	if result.Error != nil {
		return models.ChatMessage{}, fmt.Errorf("failed to create message in database: %w", result.Error)
	}
	return msg, nil
}

func (r messageRepository) GetMessages(ctx context.Context, roomId, lastSeenMsgId uint64) ([]models.ChatMessage, error) {
	var msgs []models.ChatMessage

	if lastSeenMsgId == 0 {
		result := r.db.WithContext(ctx).Where("room_id = ?", roomId).Find(&msgs)
		if result.Error != nil {
			return nil, fmt.Errorf("failed to retrieve messages for room %d: %w", roomId, result.Error)
		}
		return msgs, nil
	}

	result := r.db.WithContext(ctx).Where("room_id = ? AND id > (?)", roomId, lastSeenMsgId).Order("id asc").Find(&msgs)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to retrieve messages for room %d: %w", roomId, result.Error)
	}

	return msgs, nil
}
