package repositories

import (
	"context"
	"time"

	"github.com/i474232898/chatserver/internal/app/repositories/models"
	"gorm.io/gorm"
)

type messageRepository struct {
	db *gorm.DB
}

type MessageRepository interface {
	Create(ctx context.Context, msg *models.ChatMessage) (models.ChatMessage, error)
	GetMessages(ctx context.Context, roomId uint64, since time.Time) ([]models.ChatMessage, error)
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

func (r messageRepository) GetMessages(ctx context.Context, roomId uint64, since time.Time) ([]models.ChatMessage, error) {
	var msgs []models.ChatMessage
	result := r.db.WithContext(ctx).Where("room_id = ? AND created_at <= ?", roomId, since).Find(&msgs)
	if result.Error != nil {
		return nil, result.Error
	}

	return msgs, nil
}
