package services

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/i474232898/chatserver/api/types"
	"github.com/i474232898/chatserver/internal/app/dto"
	"github.com/i474232898/chatserver/internal/app/repositories"
	"github.com/i474232898/chatserver/internal/app/repositories/models"
)

type AuthService interface {
	Signup(ctx context.Context, user *types.SignupRequest) (*dto.SignupResponse, error)
}

type authService struct {
	userRepository repositories.UserRepository
}

func NewAuthService(userRepository repositories.UserRepository) AuthService {
	return &authService{userRepository: userRepository}
}

func (serv *authService) Signup(ctx context.Context, user *types.SignupRequest) (*dto.SignupResponse, error) {
	hashedPassword, err := serv.hashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	newUser, err := serv.userRepository.Create(ctx, &models.User{Email: string(user.Email), Password: hashedPassword})

	if err != nil {
		fmt.Errorf("failed to create user. %w", err)
		return nil, err
	}

	return &dto.SignupResponse{Email: newUser.Email, ID: newUser.ID}, nil
}

func (serv *authService) hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", fmt.Errorf("hashing failed: %w", err)
	}

	return string(hashedPassword), nil
}
