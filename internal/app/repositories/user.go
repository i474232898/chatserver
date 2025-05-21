package repositories

import (
	"context"

	// "github.com/i474232898/chatserver/internal/app/dto"
	// "github.com/i474232898/chatserver/api/types"
	// "github.com/i474232898/chatserver/internal/app/dto"
	"github.com/i474232898/chatserver/internal/app/repositories/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	result := r.db.Create(user)

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}
