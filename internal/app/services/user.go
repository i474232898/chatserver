package services

import (
	"context"
	"time"

	"github.com/i474232898/chatserver/internal/app/dto"
	"github.com/i474232898/chatserver/internal/app/repositories"
)

type UserService interface {
	Me(ctx context.Context, userId int64) (*dto.UserResponse, error)
}

type userService struct {
	userRepository repositories.UserRepository
}

func NewUserService(userRepository repositories.UserRepository) UserService {
	return &userService{userRepository: userRepository}
}

func (serv *userService) Me(ctx context.Context, userId int64) (*dto.UserResponse, error) {
	user, err := serv.userRepository.GetById(ctx, userId)
	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:        int(user.ID),
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		Username:  user.Username,
	}, nil
}
