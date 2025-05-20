package services

import (
	"context"
	"log/slog"

	"golang.org/x/crypto/bcrypt"

	"github.com/i474232898/chatserver/api/types"
	"github.com/i474232898/chatserver/internal/app/repositories"
	"github.com/i474232898/chatserver/internal/app/repositories/models"
)

type AuthService interface {
	Signup(ctx context.Context, user *types.SignupRequest) (*models.User, error)
}

type authService struct {
	userRepository repositories.UserRepository
}

func NewAuthService(userRepository repositories.UserRepository) AuthService {
	return &authService{userRepository: userRepository}
}

func (serv *authService) Signup(ctx context.Context, user *types.SignupRequest) (*models.User, error) {
	hashedPassword, err := serv.hashPassword(user.Password)
	if err != nil {
		slog.Error("Password hashing failed", "error", err.Error())
	}

	newUser, err := serv.userRepository.Create(ctx, &models.User{Email: string(user.Email), Password: hashedPassword})

	if err != nil {
		slog.Error("Failed to create user", "error", err.Error())
		return nil, err
	}

	return newUser, nil
}

func (serv *authService) hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		slog.Error("Hashing failed", "error", err.Error())
		return "", err
	}

	return string(hashedPassword), nil
}
