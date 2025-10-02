package repositories

import (
	"context"
	"fmt"

	"github.com/i474232898/chatserver/internal/app/repositories/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetById(ctx context.Context, id int64) (*models.User, error)
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
		return nil, fmt.Errorf("Unable to create user in database: %w", result.Error)
	}

	return user, nil
}

func (r *userRepository) GetById(ctx context.Context, id int64) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		return nil, fmt.Errorf("Unable to retrieve user with id %d: %w", id, err)
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, fmt.Errorf("Unable to retrieve user with email %s: %w", email, err)
	}
	return &user, nil
}
